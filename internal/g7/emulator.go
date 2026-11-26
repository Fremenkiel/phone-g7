package g7

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fremenkiel/phony-g7/internal/event"
	"github.com/fremenkiel/phony-g7/internal/platform/logger"
	"golang.org/x/net/websocket"
)

// defaultHeartbeatInterval is the interval of the documentation. The host can request a different
// interval with "HBI" in its Heartbeat response.
const defaultHeartbeatInterval = 60 * time.Second

// heartbeatGrace is the extra time that the host has to answer before the device closes the
// connection.
const heartbeatGrace = 10 * time.Second

// defaultDenyTimeout is the time of the deny screen when the host gives no "Timeout_ms". The
// documentation gives no default, so this value is an assumption.
const defaultDenyTimeout = 3 * time.Second

// The screen names of Screen.
const (
	ScreenIdle = "idle"
	ScreenDeny = "deny"
)

// Screen is what the display of the device shows.
type Screen struct {
	Name string `json:"name"`

	// Device is the reader of the screen. Empty for the idle screen.
	Device string `json:"device,omitempty"`

	// DisplayText is the text of the host. It can be HTML.
	DisplayText string `json:"display_text,omitempty"`

	// until is the time at which the device goes back to the idle screen.
	until time.Time
}

type repository interface {
	CreateIngoingEvent(ctx context.Context, payload string) error
	CreateOutgoingEvent(ctx context.Context, payload string) (int64, error)
	GetPending(ctx context.Context) ([]event.GetPendingRow, error)
	UpdateSentEvent(ctx context.Context, id int64) error
}

type Emulator struct {
	repository            repository
	Wake                  chan struct{}
	Connected             atomic.Bool
	deviceId, accessToken string

	nextTransactionID atomic.Int64

	// sendMu lets only one goroutine send the pending messages. Without it, the read loop and the
	// session loop can send the same message twice.
	sendMu sync.Mutex

	// mu guards the state of the connection below.
	mu sync.Mutex

	// subscriptions holds the event names of RegisterEvent, for example "App.*".
	subscriptions  map[string]struct{}
	onlineDecision bool

	// lastIdent is the data of the last App.CardIdent event that the device sent.
	lastIdent *CardIdentEventBody

	// screen is the screen of the last process. It is valid until its time is over.
	screen Screen

	heartbeatInterval time.Duration

	// lastHeartbeatSent is the time of the last Heartbeat request of the device.
	lastHeartbeatSent time.Time

	// lastHostHeartbeat is the time of the last Heartbeat request of the host.
	lastHostHeartbeat time.Time

	// lastHeartbeatSeen is the time of the last Heartbeat request or response of the host.
	lastHeartbeatSeen time.Time
}

func NewEmulator(r repository, deviceId, accessToken string) *Emulator {
	return &Emulator{
		repository:        r,
		Wake:              make(chan struct{}, 1),
		deviceId:          deviceId,
		accessToken:       accessToken,
		subscriptions:     map[string]struct{}{},
		heartbeatInterval: defaultHeartbeatInterval,
	}
}

// Subscriptions returns the event names that the host registered on this connection.
func (e *Emulator) Subscriptions() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	names := make([]string, 0, len(e.subscriptions))
	for name := range e.subscriptions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Screen returns the screen that the device shows now.
func (e *Emulator) Screen() Screen {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.screen.Name == "" || !time.Now().Before(e.screen.until) {
		return Screen{Name: ScreenIdle}
	}
	return e.screen
}

// SubscribedTo tells if the host registered to the event. A name that ends with "*" selects all
// events with that prefix: "App.*" selects "App.CardIdent".
func (e *Emulator) SubscribedTo(eventName CommandType) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	for name := range e.subscriptions {
		if name == string(eventName) {
			return true
		}
		if prefix, ok := strings.CutSuffix(name, "*"); ok && strings.HasPrefix(string(eventName), prefix) {
			return true
		}
	}
	return false
}

func (e *Emulator) ConnectForever(ctx context.Context, target string) {
	config, err := websocket.NewConfig(target, "http://localhost/")
	if err != nil {
		log.Printf("WebSocket configuration: %v", err)
		return
	}

	config.Header.Add("Authorization", e.accessToken)
	config.Header.Add("DeviceId", e.deviceId)
	// "new" tells the host that it must initialize the device: no registration is active.
	config.Header.Add("G7-Connection", "new")

	for ctx.Err() == nil {
		dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		ws, err := config.DialContext(dialCtx)
		cancel()
		if err == nil {
			e.resetSession()
			e.Connected.Store(true)
			log.Printf("WebSocket connected: %s", target)
			err = e.session(ctx, ws)
			e.Connected.Store(false)
			e.resetSession()
			ws.Close()
		}
		if ctx.Err() != nil {
			return
		}
		log.Printf("WebSocket unavailable: %v; retrying", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

// resetSession removes the state of a connection. A registration is valid for one connection
// only.
func (e *Emulator) resetSession() {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	e.subscriptions = map[string]struct{}{}
	e.onlineDecision = false
	e.lastIdent = nil
	e.screen = Screen{}
	e.heartbeatInterval = defaultHeartbeatInterval
	e.lastHeartbeatSent = time.Time{}
	e.lastHostHeartbeat = time.Time{}
	e.lastHeartbeatSeen = now
}

func (e *Emulator) session(ctx context.Context, ws *websocket.Conn) error {
	readErrors := make(chan error, 1)
	go func() {
		for {
			var payload json.RawMessage
			if err := websocket.JSON.Receive(ws, &payload); err != nil {
				readErrors <- err
				return
			}
			if err := e.repository.CreateIngoingEvent(ctx, string(payload)); err != nil {
				logger.LogError("emulator", "failed to log incoming message", slog.Any("error", err))
			}
			var command Command
			if err := json.Unmarshal(payload, &command); err != nil {
				logger.LogError("emulator", "unable to unmarshal payload", slog.String("payload", string(payload)))
				continue
			}

			if err := e.handleReceivedEvent(ctx, ws, &command); err != nil {
				logger.LogError(
					"emulator",
					"failed to handle message",
					slog.String("cmd", string(command.CommandType)),
					slog.Any("error", err),
				)
			}
		}
	}()

	if err := e.sendHeartbeat(ctx, ws); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := e.sendPending(ctx, ws); err != nil {
			return err
		}
		select {
		case err := <-readErrors:
			return err
		case <-e.Wake:
		case <-ticker.C:
			if err := e.heartbeatTick(ctx, ws); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// heartbeatTick sends a Heartbeat request when the interval is over and the host sent no
// Heartbeat request in that time. It returns an error when the host is silent for too long, and
// the session then closes the connection, as a device does.
func (e *Emulator) heartbeatTick(ctx context.Context, ws *websocket.Conn) error {
	e.mu.Lock()
	interval := e.heartbeatInterval
	lastSeen := e.lastHeartbeatSeen
	lastActivity := e.lastHeartbeatSent
	if e.lastHostHeartbeat.After(lastActivity) {
		lastActivity = e.lastHostHeartbeat
	}
	e.mu.Unlock()

	now := time.Now()
	if now.Sub(lastSeen) > 2*interval+heartbeatGrace {
		return fmt.Errorf("no heartbeat request or response from the host for %s", now.Sub(lastSeen).Round(time.Second))
	}
	if now.Sub(lastActivity) >= interval {
		return e.sendHeartbeat(ctx, ws)
	}
	return nil
}

// sendPending sends the queued messages in the sequence of their IDs. An event stays in the
// queue until the host registers to it: the device sends events to registered listeners only.
func (e *Emulator) sendPending(ctx context.Context, ws *websocket.Conn) error {
	e.sendMu.Lock()
	defer e.sendMu.Unlock()

	events, err := e.repository.GetPending(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	for i := range events {
		event := events[i]

		// A payload of the raw trigger can be a JSON object that is not a command. Such a payload
		// has no message type and the device sends it as it is.
		var command Command
		_ = json.Unmarshal([]byte(event.Payload), &command)

		if command.MessageType == MessageTypeEvent && !e.SubscribedTo(command.CommandType) {
			continue
		}

		_ = ws.SetWriteDeadline(time.Time{})

		if err := websocket.JSON.Send(ws, json.RawMessage(event.Payload)); err != nil {
			return fmt.Errorf("send event %d: %w", event.ID, err)
		}
		if err := e.repository.UpdateSentEvent(ctx, event.ID); err != nil {
			return err
		}

		if command.MessageType == MessageTypeEvent && command.CommandType == CmdAppCardIdent {
			if ident, err := DecodeData[CardIdentEventBody](command); err == nil {
				e.mu.Lock()
				e.lastIdent = &ident
				e.mu.Unlock()
			}
		}
	}

	return nil
}

// sendHeartbeat sends a Heartbeat request. The data object of the request is empty.
func (e *Emulator) sendHeartbeat(ctx context.Context, ws *websocket.Conn) error {
	command, err := NewCommand(CmdHeartbeat, MessageTypeRequest, nil)
	if err != nil {
		return err
	}
	command.TransactionID = new(e.nextTransactionID.Add(1))

	e.mu.Lock()
	e.lastHeartbeatSent = time.Now()
	e.mu.Unlock()

	return e.sendCommand(ctx, ws, &command)
}

func (e *Emulator) handleReceivedEvent(ctx context.Context, ws *websocket.Conn, payload *Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	switch payload.MessageType {
	case MessageTypeRequest:
		return e.handleRequest(ctx, ws, payload)
	case MessageTypeResponse:
		return e.handleResponse(payload)
	case MessageTypeEvent:
		return nil
	}

	return nil
}

// handleResponse reads the Heartbeat response of the host. The host can request a different
// heartbeat interval with "HBI".
func (e *Emulator) handleResponse(payload *Command) error {
	if payload.CommandType != CmdHeartbeat {
		return nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.lastHeartbeatSeen = time.Now()

	body, err := DecodeData[heartbeatInterval](*payload)
	if err != nil {
		return err
	}
	if body.HeartbeatInterval > 0 {
		e.heartbeatInterval = time.Duration(body.HeartbeatInterval) * time.Second
	}
	return nil
}

// heartbeatInterval reads only "HBI" of a Heartbeat response. The host can write "RT" in a time
// format that HeartbeatBody does not read.
type heartbeatInterval struct {
	HeartbeatInterval int `json:"HBI"`
}

func (e *Emulator) handleRequest(ctx context.Context, ws *websocket.Conn, payload *Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	switch payload.CommandType {
	case CmdHeartbeat:
		now := time.Now()
		e.mu.Lock()
		e.lastHostHeartbeat = now
		e.lastHeartbeatSeen = now
		interval := e.heartbeatInterval
		e.mu.Unlock()

		return e.respond(ctx, ws, payload, StateOK, HeartbeatBody{
			HeartbeatInterval: int(interval / time.Second),
			CurrentTime:       now.UTC(),
		})

	case CmdSystemSetHostInfo:
		body, err := DecodeData[SetHostInfoRequestBody](*payload)
		if err != nil {
			return e.respond(ctx, ws, payload, StateParameterError, nil)
		}
		e.mu.Lock()
		e.onlineDecision = body.OnlineDecision != nil && *body.OnlineDecision
		e.mu.Unlock()
		return e.respond(ctx, ws, payload, StateOK, nil)

	case CmdGetDeviceInfo:
		return e.respond(ctx, ws, payload, StateOK, e.deviceInfo())

	case CmdRegisterEvent:
		body, err := DecodeData[RegisterEventRequestBody](*payload)
		if err != nil || body.Event == "" {
			return e.respond(ctx, ws, payload, StateParameterError, nil)
		}
		if err := e.respond(ctx, ws, payload, StateOK, nil); err != nil {
			return err
		}
		// The registration starts after the response, or a queued event can go out before the
		// response.
		e.mu.Lock()
		e.subscriptions[body.Event] = struct{}{}
		e.mu.Unlock()
		return e.sendPending(ctx, ws)

	case CmdUnregisterEvent:
		body, err := DecodeData[RegisterEventRequestBody](*payload)
		if err != nil || body.Event == "" {
			return e.respond(ctx, ws, payload, StateParameterError, nil)
		}
		e.mu.Lock()
		delete(e.subscriptions, body.Event)
		e.mu.Unlock()
		return e.respond(ctx, ws, payload, StateOK, nil)

	case CmdAppStartUnlockProcess:
		body, err := DecodeData[StartUnlockProcessRequestBody](*payload)
		if err != nil {
			return e.respond(ctx, ws, payload, StateParameterError, nil)
		}
		if err := e.respond(ctx, ws, payload, StateOK, nil); err != nil {
			return err
		}
		return e.sendUnlockProcessResult(ctx, ws, body.Device)

	case CmdAppStartDenyProcess:
		return e.handleStartDenyProcess(ctx, ws, payload)
	}

	// The emulator does not have this command. The protocol gives state 255 for such a request.
	return e.respond(ctx, ws, payload, StateUnknownCommand, nil)
}

// handleStartDenyProcess does the deny process of App.StartDenyProcess: the device shows the deny
// screen, then goes to idle. The identification that waited for a decision is complete, so a
// later unlock does not get its ID. The device sends no event for a deny: the documentation has
// no result event, and its multiple door example shows only the response.
func (e *Emulator) handleStartDenyProcess(ctx context.Context, ws *websocket.Conn, payload *Command) error {
	body, err := DecodeData[StartDenyProcessRequestBody](*payload)
	if err != nil || (body.TimeoutMs != nil && *body.TimeoutMs < 0) {
		return e.respond(ctx, ws, payload, StateParameterError, nil)
	}

	screen := Screen{
		Name:        ScreenDeny,
		Device:      DeviceInternal,
		DisplayText: "Access denied",
		until:       time.Now().Add(defaultDenyTimeout),
	}
	if body.Device != nil {
		screen.Device = *body.Device
	}
	if body.DisplayText != nil {
		screen.DisplayText = *body.DisplayText
	}
	if body.TimeoutMs != nil {
		screen.until = time.Now().Add(time.Duration(*body.TimeoutMs) * time.Millisecond)
	}

	e.mu.Lock()
	e.lastIdent = nil
	e.screen = screen
	e.mu.Unlock()

	return e.respond(ctx, ws, payload, StateOK, nil)
}

// respond sends the response to a request. The response has the command name, the transaction ID
// and the global message ID of the request.
func (e *Emulator) respond(ctx context.Context, ws *websocket.Conn, request *Command, state State, body any) error {
	command, err := NewCommand(request.CommandType, MessageTypeResponse, body)
	if err != nil {
		return err
	}
	command.TransactionID = request.TransactionID
	command.GlobalMessageID = request.GlobalMessageID
	command.State = new(state)
	return e.sendCommand(ctx, ws, &command)
}

// sendUnlockProcessResult sends the event App.UnlockProcessResult: the unlock process is complete
// and the access was granted. A device sends it after the unlocking time. The emulator sends it
// immediately.
func (e *Emulator) sendUnlockProcessResult(ctx context.Context, ws *websocket.Conn, device *string) error {
	e.mu.Lock()
	ident := e.lastIdent
	e.lastIdent = nil
	e.screen = Screen{}
	e.mu.Unlock()

	result := UnlockProcessResultEventBody{
		AccessGranted: true,
		Feedback:      true,
		Count:         new(1),
		Device:        device,
	}

	if ident != nil {
		if result.Device == nil {
			result.Device = ident.Device
		}
		result.CardTypeID = ident.CardTypeID

		cardType := ""
		if ident.CardType != nil {
			cardType = *ident.CardType
		}
		switch {
		case cardType == CardTypePIN && len(ident.Segments) > 0:
			result.PIN = ident.Segments[0].Data
		case cardType == CardTypeBarcode && len(ident.Segments) > 0:
			result.Barcode = ident.Segments[0].BarcodeID
		default:
			result.UID = ident.UID
		}
	}
	if result.Device == nil {
		result.Device = new(DeviceInternal)
	}

	command, err := NewCommand(CmdAppUnlockProcessResult, MessageTypeEvent, result)
	if err != nil {
		return err
	}
	return e.sendCommand(ctx, ws, &command)
}

// deviceInfo is the response of GetDeviceInfo. The device ID has the form
// "<ArticleNumber>-<SerialNumber>".
func (e *Emulator) deviceInfo() GetDeviceInfoResponseBody {
	info := GetDeviceInfoResponseBody{
		DeviceModel:     new("GT7 (phony-g7 emulator)"),
		DeviceTypeID:    new(113),
		DeviceName:      new("phony-g7"),
		FirmwareVersion: new("3.12.0"),
		ActiveApp: &DeviceActiveApp{
			Name:        new("access"),
			DisplayName: new("G7 Access"),
		},
	}

	if i := strings.LastIndex(e.deviceId, "-"); i > 0 {
		info.ArticleNumber = new(e.deviceId[:i])
		info.SerialNumber = new(e.deviceId[i+1:])
	} else {
		info.SerialNumber = new(e.deviceId)
	}
	return info
}

func (e *Emulator) sendCommand(ctx context.Context, ws *websocket.Conn, command *Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	payload, err := json.Marshal(command)
	if err != nil {
		return err
	}

	if _, err := e.repository.CreateOutgoingEvent(ctx, string(payload)); err != nil {
		return err
	}

	if err := e.sendPending(ctx, ws); err != nil {
		return err
	}

	return nil
}
