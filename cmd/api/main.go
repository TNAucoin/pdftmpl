package main

import (
	"fmt"
	"github.com/tnaucoin/pdftmpl/cmd/api/router"
	"github.com/tnaucoin/pdftmpl/config"
	"github.com/tnaucoin/pdftmpl/internal/pdfClient"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	config *config.Conf
	logger *slog.Logger
}

func main() {
	cfg := config.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	gotenClient := pdfClient.New(cfg.Goten, logger)

	r := router.New(logger, gotenClient)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		IdleTimeout:  cfg.Server.TimeoutIdle,
		WriteTimeout: cfg.Server.TimeoutWrite,
		ReadTimeout:  cfg.Server.TimeoutRead,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "port", cfg.Server.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
