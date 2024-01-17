package main

import (
	"context"
	"errors"
	"fmt"
	apiRouter "github.com/tnaucoin/pdftmpl/cmd/api/router"
	"github.com/tnaucoin/pdftmpl/config"
	"github.com/tnaucoin/pdftmpl/internal/pdfClient"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main is the entry point of the application
func main() {
	cfg := config.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	r := router(cfg, logger)
	srv := configureServer(cfg, r, logger)
	runServer(srv, logger)
}

// configureServer creates and configures an HTTP server.
// The server listens on the specified port from the configuration.
// The router is used to handle incoming requests.
// The logger is used for error logging.
// The server's idle timeout, write timeout, and read timeout are set using the values from the configuration.
// The server's error log is set to a logger with a log level of LevelError based on the provided logger.
// Returns the configured HTTP server.
func configureServer(cfg *config.Conf, router http.Handler, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		IdleTimeout:  cfg.Server.TimeoutIdle,
		WriteTimeout: cfg.Server.TimeoutWrite,
		ReadTimeout:  cfg.Server.TimeoutRead,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}

// router creates the application routes for the api
// it also initializes a gotenClient for PDF generation
func router(cfg *config.Conf, logger *slog.Logger) http.Handler {
	gotenClient := pdfClient.New(cfg.Goten, logger)
	return apiRouter.New(logger, gotenClient)
}

// runServer server run context
func runServer(srv *http.Server, logger *slog.Logger) {
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
