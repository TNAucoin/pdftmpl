package generate

import (
	"bytes"
	"context"
	"errors"
	"github.com/tnaucoin/pdftmpl/internal/weasyPrintClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"github.com/tnaucoin/pdftmpl/utils"
	"log/slog"
	"strconv"
)

type InvoiceGenerateRequest struct {
	Body struct {
		InvoiceID        uint64 `json:"invoice-id,required" doc:"Invoice ID number" example:"1024"`
		RecipientName    string `json:"recipient-name,required" maxLength:"50" doc:"Name of invoice recipient" example:"John Doe"`
		RecipientAddress string `json:"recipient-address,required" doc:"Address of invoice recipient" example:"123 Somewhere Ln."`
		RecipientEmail   string `json:"recipient-email,required" doc:"Email of invoice recipient" example:"example@gmail.com" pattern:"^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"`
		Logo             string `json:"logo,required" enum:"logo.png,logo2.png" doc:"Logo to include on invoice"`
	}
}

type InvoiceGenerateResponse struct {
}

type Handler struct {
	logger  *slog.Logger
	wpc     *weasyPrintClient.WeasyPrintClient
	outPath string
}

func New(logger *slog.Logger, wpc *weasyPrintClient.WeasyPrintClient, outPath string) *Handler {
	return &Handler{
		logger:  logger,
		wpc:     wpc,
		outPath: outPath,
	}
}

func (h *Handler) GenerateInvoicePDFHandler(_ context.Context, input *InvoiceGenerateRequest) (*InvoiceGenerateResponse, error) {
	file, err := h.createInvoicePDF(input)
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	err = utils.WriteToFile(file, h.outPath, "test-pdf.pdf")
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	return nil, nil
}

func (h *Handler) createInvoicePDF(input *InvoiceGenerateRequest) (*bytes.Buffer, error) {
	doc, err := utils.RenderTemplate(templates.Hello(
		strconv.FormatUint(input.Body.InvoiceID, 10),
		input.Body.RecipientName,
		input.Body.RecipientAddress,
		input.Body.RecipientEmail,
		input.Body.Logo,
	))
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	var pdf bytes.Buffer
	err = h.wpc.GeneratePDF(doc, &pdf)
	if err != nil {
		if errors.Is(err, weasyPrintClient.ErrWeasyPrintGeneration) {
			return nil, err
		}
		// Log errors not caused by weasyPrintClient
		h.logger.Error(err.Error())
		return nil, err
	}

	return &pdf, nil
}
