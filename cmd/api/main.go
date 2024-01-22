package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/tnaucoin/pdftmpl/cmd/api/router"
	"github.com/tnaucoin/pdftmpl/config"
	"github.com/tnaucoin/pdftmpl/internal/weasyPrintClient"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// main initializes and starts the PDF Generator API server.
// It performs the following steps:
// - Creates a new configuration using config.New().
// - Creates a new logger with debug level using slog.New().
// - Creates a new WeasyPrint client using the logger with weasyPrintClient.New().
// - Creates a new Application instance with the logger, configuration, and WeasyPrint client.
// - Configures the HTTP server using the configureServer() function with the Application instance.
// - Starts the server by calling the runServer() function with the configured server instance.
// This function does not return any value.
func main() {
	cfg := config.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: setLogLevel(cfg.Server.LogLevel),
	}))
	wpc := weasyPrintClient.New(logger)
	mux := router.New(logger, wpc, cfg)
	srv := configureServer(cfg, logger, mux)
	runServer(srv)
}

// configureServer creates and configures an HTTP server based on the provided Application struct.
// The server is configured with the following properties:
// - Addr: The server address, composed of the configured server port from the Application's config.
// - Handler: The router returned from the Application's Router method.
// - IdleTimeout: The idle timeout value from the Application's config.
// - WriteTimeout: The write timeout value from the Application's config.
// - ReadTimeout: The read timeout value from the Application's config.
// - ErrorLog: A logger configured with the Application's logger and the logging level set to Error.
// This function returns the configured http.Server instance.
func configureServer(cfg *config.Conf, logger *slog.Logger, mux *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		IdleTimeout:  cfg.Server.TimeoutIdle,
		WriteTimeout: cfg.Server.TimeoutWrite,
		ReadTimeout:  cfg.Server.TimeoutRead,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}

// runServer server run context
func runServer(srv *http.Server) {
	// server run context
	serverCtx, serverStopCtx := context.WithCancelCause(context.Background())
	// listen for syscall signals for process to interrupt/quit
	setupGracefulShutdown(serverCtx, serverStopCtx, srv)

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	// wait for server context to be stopped
	<-serverCtx.Done()
}

// setupGracefulShutdown listens for syscall signals and triggers graceful shutdown of the HTTP server.
// The function takes the following parameters:
// - serverCtx: Context used to manage the server's lifecycle.
// - serverStopCtx: A function to stop the server's context.
// - srv: Pointer to the HTTP server instance.
// - logger: Pointer to the logger instance.
func setupGracefulShutdown(serverCtx context.Context, serverStopCtx context.CancelCauseFunc, srv *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		// shutdown grace period
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out, forcing exit.")
			}
		}()
		// trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx(nil)
	}()
}

// setLogLevel converts the given log level string to a slog.Level.
// It performs a case-insensitive comparison and returns the corresponding slog.Level value.
// If the log level is not recognized, it defaults to slog.LevelInfo.
func setLogLevel(logLevel string) slog.Level {
	switch strings.ToLower(logLevel) {
	case "error":
		return slog.LevelError
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "debug":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
