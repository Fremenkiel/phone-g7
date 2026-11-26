package g7

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fremenkiel/phony-g7/internal/event"
	"golang.org/x/net/websocket"
)

// memoryRepository is the events table in memory.
type memoryRepository struct {
	mu     sync.Mutex
	nextID int64
	queued []event.GetPendingRow
}

func (m *memoryRepository) CreateIngoingEvent(context.Context, string) error { return nil }

func (m *memoryRepository) CreateOutgoingEvent(_ context.Context, payload string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	m.queued = append(m.queued, event.GetPendingRow{ID: m.nextID, Payload: payload})
	return m.nextID, nil
}

func (m *memoryRepository) GetPending(context.Context) ([]event.GetPendingRow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]event.GetPendingRow{}, m.queued...), nil
}

func (m *memoryRepository) UpdateSentEvent(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.queued {
		if m.queued[i].ID == id {
			m.queued = append(m.queued[:i], m.queued[i+1:]...)
			break
		}
	}
	return nil
}

// host is the management software: the WebSocket server that the device connects to.
type host struct {
	t        *testing.T
	messages chan Command
	conn     chan *websocket.Conn
	header   chan http.Header
}

func newHost(t *testing.T) (*host, string) {
	h := &host{
		t:        t,
		messages: make(chan Command, 32),
		conn:     make(chan *websocket.Conn, 1),
		header:   make(chan http.Header, 1),
	}
	server := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		h.header <- ws.Request().Header
		h.conn <- ws
		for {
			var command Command
			if err := websocket.JSON.Receive(ws, &command); err != nil {
				return
			}
			h.messages <- command
		}
	}))
	t.Cleanup(server.Close)
	return h, "ws" + strings.TrimPrefix(server.URL, "http")
}

func (h *host) next() Command {
	h.t.Helper()
	select {
	case command := <-h.messages:
		return command
	case <-time.After(3 * time.Second):
		h.t.Fatal("the device sent no message")
		return Command{}
	}
}

func (h *host) silent(d time.Duration) {
	h.t.Helper()
	select {
	case command := <-h.messages:
		h.t.Fatalf("the device sent %s %s, want no message", command.CommandType, command.MessageType)
	case <-time.After(d):
	}
}

func (h *host) send(ws *websocket.Conn, cmd CommandType, mt MessageType, tid int64, body any) {
	h.t.Helper()
	command, err := NewCommand(cmd, mt, body)
	if err != nil {
		h.t.Fatal(err)
	}
	command.TransactionID = &tid
	if mt == MessageTypeResponse {
		command.State = new(StateOK)
	}
	if err := websocket.JSON.Send(ws, command); err != nil {
		h.t.Fatal(err)
	}
}

func (h *host) wantResponse(cmd CommandType, tid int64, state State) Command {
	h.t.Helper()
	command := h.next()
	if command.CommandType != cmd || command.MessageType != MessageTypeResponse {
		h.t.Fatalf("got %s %s, want %s Rsp", command.CommandType, command.MessageType, cmd)
	}
	if command.TransactionID == nil || *command.TransactionID != tid {
		h.t.Fatalf("%s: got TID %v, want %d", cmd, command.TransactionID, tid)
	}
	if command.State == nil || *command.State != state {
		h.t.Fatalf("%s: got state %v, want %d", cmd, command.State, state)
	}
	return command
}

func queue(t *testing.T, e *Emulator, repo *memoryRepository, command Command, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(command)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateOutgoingEvent(context.Background(), string(payload)); err != nil {
		t.Fatal(err)
	}
	select {
	case e.Wake <- struct{}{}:
	default:
	}
}

func TestSession(t *testing.T) {
	h, target := newHost(t)
	repo := &memoryRepository{}
	emulator := NewEmulator(repo, "GT7-0001", "token")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go emulator.ConnectForever(ctx, target)

	header := <-h.header
	if header.Get("Authorization") != "token" || header.Get("DeviceId") != "GT7-0001" || header.Get("G7-Connection") != "new" {
		t.Fatalf("unexpected headers: %v", header)
	}
	ws := <-h.conn

	// The device sends a Heartbeat request with an empty data object at the start.
	heartbeat := h.next()
	if heartbeat.CommandType != CmdHeartbeat || heartbeat.MessageType != MessageTypeRequest {
		t.Fatalf("got %s %s, want Heartbeat Req", heartbeat.CommandType, heartbeat.MessageType)
	}
	if heartbeat.TransactionID == nil || string(heartbeat.Data) != "{}" || heartbeat.State != nil {
		t.Fatalf("unexpected heartbeat request: %+v data %s", heartbeat, heartbeat.Data)
	}
	h.send(ws, CmdHeartbeat, MessageTypeResponse, *heartbeat.TransactionID, map[string]any{"HBI": 1, "RT": "2026-01-01T00:00:00"})

	// The host requested an interval of 1 second.
	heartbeat = h.next()
	if heartbeat.CommandType != CmdHeartbeat || heartbeat.MessageType != MessageTypeRequest {
		t.Fatalf("got %s %s, want Heartbeat Req", heartbeat.CommandType, heartbeat.MessageType)
	}
	h.send(ws, CmdHeartbeat, MessageTypeResponse, *heartbeat.TransactionID, map[string]any{"HBI": 60})

	// The device answers a Heartbeat request of the host.
	h.send(ws, CmdHeartbeat, MessageTypeRequest, 7, nil)
	response := h.wantResponse(CmdHeartbeat, 7, StateOK)
	body, err := DecodeData[HeartbeatBody](response)
	if err != nil || body.HeartbeatInterval != 60 || body.CurrentTime.IsZero() {
		t.Fatalf("unexpected heartbeat response: %s (%v)", response.Data, err)
	}

	// No host is registered: the event stays in the queue.
	command, err := NewCardReadEvent("04FC8862B93F80")
	queue(t, emulator, repo, command, err)
	h.silent(1500 * time.Millisecond)

	h.send(ws, CmdSystemSetHostInfo, MessageTypeRequest, 8, SetHostInfoRequestBody{OnlineDecision: new(true)})
	h.wantResponse(CmdSystemSetHostInfo, 8, StateOK)

	// The response to RegisterEvent comes before the queued event.
	h.send(ws, CmdRegisterEvent, MessageTypeRequest, 9, RegisterEventRequestBody{Event: "App.*"})
	h.wantResponse(CmdRegisterEvent, 9, StateOK)

	ident := h.next()
	if ident.CommandType != CmdAppCardIdent || ident.MessageType != MessageTypeEvent || ident.TransactionID != nil {
		t.Fatalf("got %s %s, want App.CardIdent Evt", ident.CommandType, ident.MessageType)
	}

	h.send(ws, CmdAppStartUnlockProcess, MessageTypeRequest, 10, StartUnlockProcessRequestBody{})
	unlock := h.wantResponse(CmdAppStartUnlockProcess, 10, StateOK)
	if string(unlock.Data) != "{}" {
		t.Fatalf("unlock response data: %s", unlock.Data)
	}
	result := h.next()
	if result.CommandType != CmdAppUnlockProcessResult || result.MessageType != MessageTypeEvent {
		t.Fatalf("got %s %s, want App.UnlockProcessResult Evt", result.CommandType, result.MessageType)
	}
	resultBody, err := DecodeData[UnlockProcessResultEventBody](result)
	if err != nil || !resultBody.AccessGranted || resultBody.UID == nil || *resultBody.UID != "04FC8862B93F80" || *resultBody.Device != DeviceInternal {
		t.Fatalf("unexpected unlock result: %s (%v)", result.Data, err)
	}

	// PIN, then deny: a response and no event.
	command, err = NewPinEvent("1234")
	queue(t, emulator, repo, command, err)
	if ident = h.next(); ident.CommandType != CmdAppCardIdent {
		t.Fatalf("got %s, want App.CardIdent", ident.CommandType)
	}
	h.send(ws, CmdAppStartDenyProcess, MessageTypeRequest, 11, StartDenyProcessRequestBody{
		DisplayText: new("Denied!"),
		TimeoutMs:   new(60000),
	})
	deny := h.wantResponse(CmdAppStartDenyProcess, 11, StateOK)
	if string(deny.Data) != "{}" {
		t.Fatalf("deny response data: %s", deny.Data)
	}
	if screen := emulator.Screen(); screen.Name != ScreenDeny || screen.DisplayText != "Denied!" || screen.Device != DeviceInternal {
		t.Fatalf("unexpected screen after deny: %+v", screen)
	}
	h.silent(500 * time.Millisecond)

	// The deny completed the identification: a later unlock does not get the PIN.
	h.send(ws, CmdAppStartUnlockProcess, MessageTypeRequest, 20, StartUnlockProcessRequestBody{})
	h.wantResponse(CmdAppStartUnlockProcess, 20, StateOK)
	resultBody, err = DecodeData[UnlockProcessResultEventBody](h.next())
	if err != nil || resultBody.PIN != nil || resultBody.UID != nil {
		t.Fatalf("unexpected unlock result after deny: %+v (%v)", resultBody, err)
	}
	if screen := emulator.Screen(); screen.Name != ScreenIdle {
		t.Fatalf("unexpected screen after unlock: %+v", screen)
	}

	// A deny with a wrong parameter type gets state 2.
	h.send(ws, CmdAppStartDenyProcess, MessageTypeRequest, 21, map[string]any{"Timeout_ms": "long"})
	h.wantResponse(CmdAppStartDenyProcess, 21, StateParameterError)

	// Barcode, then unlock: the result has the barcode.
	command, err = NewBarcodeEvent("123456789")
	queue(t, emulator, repo, command, err)
	ident = h.next()
	identBody, err := DecodeData[CardIdentEventBody](ident)
	if err != nil || identBody.UID == nil || *identBody.UID != "075BCD15" || *identBody.Device != "BARCODE1" {
		t.Fatalf("unexpected barcode event: %s (%v)", ident.Data, err)
	}
	h.send(ws, CmdAppStartUnlockProcess, MessageTypeRequest, 12, StartUnlockProcessRequestBody{})
	h.wantResponse(CmdAppStartUnlockProcess, 12, StateOK)
	resultBody, err = DecodeData[UnlockProcessResultEventBody](h.next())
	if err != nil || resultBody.Barcode == nil || *resultBody.Barcode != "123456789" || resultBody.UID != nil {
		t.Fatalf("unexpected unlock result: %+v (%v)", resultBody, err)
	}

	// A command that the emulator does not have gets state 255.
	h.send(ws, CmdAppGetCameraImage, MessageTypeRequest, 13, nil)
	h.wantResponse(CmdAppGetCameraImage, 13, StateUnknownCommand)

	// After UnregisterEvent, events stay in the queue again.
	h.send(ws, CmdUnregisterEvent, MessageTypeRequest, 14, RegisterEventRequestBody{Event: "App.*"})
	h.wantResponse(CmdUnregisterEvent, 14, StateOK)
	command, err = NewCardReadEvent("AABBCCDD")
	queue(t, emulator, repo, command, err)
	h.silent(1500 * time.Millisecond)
}

func TestIdentEvents(t *testing.T) {
	if _, err := NewCardReadEvent("  "); err != ErrIdentEmpty {
		t.Fatalf("got %v, want ErrIdentEmpty", err)
	}
	if _, err := NewPinEvent("12a4"); err != ErrPINNotDigits {
		t.Fatalf("got %v, want ErrPINNotDigits", err)
	}

	command, err := NewPinEvent("1234")
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(command)
	want := `{"Cmd":"App.CardIdent","MT":"Evt","Data":{"UID":"0","RfStandard":"UNKNOWN","CardType":"PIN","Segments":[{"SegmentType":"PIN_DATA","Data":"1234","PINType":"IDENTIFICATION"}],"Language":"en","Device":"INTERNAL"}}`
	if string(payload) != want {
		t.Fatalf("got  %s\nwant %s", payload, want)
	}

	command, err = NewCardReadEvent("04FC8862B93F80")
	if err != nil {
		t.Fatal(err)
	}
	payload, _ = json.Marshal(command)
	if !strings.Contains(string(payload), `"Segments":[]`) {
		t.Fatalf("a card read has an empty segment list: %s", payload)
	}
}
