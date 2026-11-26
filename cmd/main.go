package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fremenkiel/phony-g7/internal/action"
	"github.com/fremenkiel/phony-g7/internal/admin"
	"github.com/fremenkiel/phony-g7/internal/event"
	"github.com/fremenkiel/phony-g7/internal/g7"
	"github.com/fremenkiel/phony-g7/internal/platform/logger"
	"github.com/fremenkiel/phony-g7/internal/status"
	"github.com/fremenkiel/phony-g7/internal/ui"
	uifs "github.com/fremenkiel/phony-g7/ui"

	"database/sql"

	_ "modernc.org/sqlite"
)

func main() {
	logger.NewLogger()

	// The read loop and the send loop of the emulator write to the database at the same time.
	// Without busy_timeout, SQLite does not wait for the other writer and gives SQLITE_BUSY.
	db, err := sql.Open("sqlite", "file:g7.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		panic(err)
	}
	// SQLite permits one writer at a time, so one connection is sufficient.
	db.SetMaxOpenConns(1)

	if len(os.Args) < 5 {
		logger.LogError("main", "add --device-id and --token")
		return
	}

	var deviceId, accessToken, target string
	for i := range os.Args {
		switch os.Args[i] {
		case "--device-id":
			deviceId = os.Args[i + 1]
		case "--token":
			accessToken = os.Args[i + 1]
		case "--target":
			target = os.Args[i + 1]
		}
	}

	eventRepo := event.New(db)
	eventRepo.RemoveEvents(context.Background())

	emulator := g7.NewEmulator(eventRepo, deviceId, accessToken)

	if len(target) == 0 {
		target = "ws://localhost:5010/ws/19"
	}
	go emulator.ConnectForever(context.Background(), target)

	eventService := event.NewService(eventRepo)
	actionService := action.NewService(eventRepo)

	renderer := ui.NewRenderer()

	adminHandler := admin.NewHandler(renderer)
	actionHandler := action.NewHandler(emulator, actionService)
	eventHandler := event.NewHandler(eventService)
	statusHandler := status.NewHandler(emulator)

	mux := http.NewServeMux()

	staticFiles, err := fs.Sub(uifs.Files, "static")
  if err != nil {
          logger.LogFatal("main", "unable to load static files", slog.Any("error", err))
  }
  mux.Handle("/scripts/", http.FileServer(http.FS(staticFiles)))
  mux.Handle("/styles/", http.FileServer(http.FS(staticFiles)))
  mux.Handle("/images/", http.FileServer(http.FS(staticFiles)))

	action.NewRouter(
		mux,
		actionHandler,
		)

	admin.NewRouter(
		mux,
		adminHandler,
		)

	event.NewRouter(
		mux,
		eventHandler,
		)

	status.NewRouter(
		mux,
		statusHandler,
		)

	port := "8080"
	address := fmt.Sprintf(":%s", port)

	logger.LogInfo(
		"main",
		"server started",
		slog.String("address", "127.0.0.1"),
		slog.String("port", port),
		)

	server := http.Server{
		Addr: address,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.LogInfo(
				"main",
				"http listen and serve error",
				slog.Any("error", err),
				)
		}
	}()

	<-ctx.Done()
	startShutdown := time.Now()
	logger.LogInfo(
		"main",
		"shutdown signel recived",
		slog.Int64("recived_at", startShutdown.Unix()),
		)

	shotdownCtx, shotdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shotdownRelease()

	if err := server.Shutdown(shotdownCtx); err != nil {
		logger.LogInfo(
			"main",
			"unable to gracefully shutdown http server",
			slog.Any("error", err),
			)
	}

	logger.LogInfo(
		"main",
		"shutdown complete",
		slog.Int64("duration_ms", time.Since(startShutdown).Milliseconds()),
		)
}
