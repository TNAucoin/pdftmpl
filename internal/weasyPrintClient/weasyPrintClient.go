package weasyPrintClient

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

var (
	ErrWeasyPrintGeneration = errors.New("PDF generation failed to execute")
)

type WeasyPrintClient struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *WeasyPrintClient {
	return &WeasyPrintClient{
		logger: logger,
	}
}

func (wp *WeasyPrintClient) GeneratePDF(in io.Reader, out io.Writer) error {
	startTime := time.Now()
	c := exec.Command("weasyprint", "-", "-")
	c.Stdin = in
	stdout, err := c.StdoutPipe()
	if err != nil {
		wp.logger.Error(err.Error())
		return err
	}
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		wp.logger.Error(err.Error())
		return ErrWeasyPrintGeneration
	}
	if _, err := io.Copy(out, stdout); err != nil {
		wp.logger.Error(err.Error())
		return ErrWeasyPrintGeneration
	}
	if err := c.Wait(); err != nil {
		wp.logger.Error(err.Error())
		return ErrWeasyPrintGeneration
	}
	wp.logger.Debug("weasyPrintClient", "method", "GeneratePDF", "duration", time.Since(startTime))
	return nil
}
