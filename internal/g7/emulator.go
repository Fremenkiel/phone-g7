package g7

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/fremenkiel/phony-g7/internal/event"
	"github.com/fremenkiel/phony-g7/internal/platform/logger"
	"golang.org/x/net/websocket"
)

type repository interface {
	CreateIngoingEvent(ctx context.Context, payload string) error
	CreateOutgoingEvent(ctx context.Context, payload string) (int64, error)
	GetPending(ctx context.Context) ([]event.GetPendingRow, error)
	UpdateSentEvent(ctx context.Context, id int64) error
}

type Emulator struct {
	repository						repository
	Wake									chan struct{}
	Connected							atomic.Bool
	lastHeartbeatReceived time.Time
	deviceId, accessToken string
}

func NewEmulator(r repository, deviceId, accessToken string) *Emulator {
	return &Emulator{
		repository: r,
		Wake: make(chan struct{}, 1),
		deviceId: deviceId,
		accessToken: accessToken,
	}
}

func (e *Emulator) ConnectForever(ctx context.Context, target string) {
	config, err := websocket.NewConfig(target, "http://localhost/")
	if err != nil {
		log.Printf("WebSocket configuration: %v", err)
		return
	}

	config.Header.Add("Authorization", e.accessToken)
	config.Header.Add("DeviceId", e.deviceId)

	for ctx.Err() == nil {
		dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		ws, err := config.DialContext(dialCtx)
		cancel()
		if err == nil {
			e.Connected.Store(true)
			log.Printf("WebSocket connected: %s", target)
			err = e.session(ctx, ws)
			e.Connected.Store(false)
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
			}

			e.handleReceivedEvent(ctx, ws, &command)
		}
	}()

	if err := e.sendHeartbeat(ctx, ws); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	heartbeat := time.NewTicker(60 * time.Second)
	defer heartbeat.Stop()

	for {
		if err := e.sendPending(ctx, ws); err != nil {
			return err
		}
		select {
		case err := <-readErrors:
			return err
		case <-e.Wake:
		case <-ticker.C:
		case <-heartbeat.C:
			if err := e.sendHeartbeat(ctx, ws); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *Emulator) sendPending(ctx context.Context, ws *websocket.Conn) error {
	events, err := e.repository.GetPending(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	for i := range events {
		event := events[i]

		_ = ws.SetWriteDeadline(time.Time{})

		if err := websocket.JSON.Send(ws, json.RawMessage(event.Payload)); err != nil {
			return fmt.Errorf("send event %d: %w", event.ID, err)
		}
		if err := e.repository.UpdateSentEvent(ctx, event.ID); err != nil {
			return err
		}
	}

	return nil
}

func (e *Emulator) sendHeartbeat(ctx context.Context, ws *websocket.Conn) error {
	transactionId := int64(rand.IntN(9999))
	command, err := NewCommand(CmdHeartbeat, MessageTypeRequest, HeartbeatBody{})
	if err != nil {
		return err
	}
	command.TransactionID = &transactionId

	e.sendCommand(ctx, ws, &command)

	return nil
}

func (e *Emulator) handleReceivedEvent(ctx context.Context, ws *websocket.Conn, payload *Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	switch payload.MessageType {
	case MessageTypeRequest:
		return e.handleRequest(ctx, ws, payload)
	case MessageTypeResponse:
		return nil
	case MessageTypeEvent:
		return nil
	}


	return nil
}

func (e *Emulator) handleRequest(ctx context.Context, ws *websocket.Conn, payload *Command) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	switch payload.CommandType {
	case CmdAppStartUnlockProcess:
		if cmd, err := NewCommand(CmdAppStartUnlockProcess, MessageTypeResponse, StartUnlockProcessRequestBody{}); err != nil {
			return err
		} else {
			cmd.TransactionID = payload.TransactionID
			cmd.GlobalMessageID = payload.GlobalMessageID
			cmd.State = new(StateOK)
			return e.sendCommand(ctx, ws, &cmd)
		}
	}


	return nil
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
