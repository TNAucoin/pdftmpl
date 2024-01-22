package ticket

import (
	"bytes"
	"context"
	"github.com/tnaucoin/pdftmpl/internal/weasyPrintClient"
	"github.com/tnaucoin/pdftmpl/templates"
	"github.com/tnaucoin/pdftmpl/utils"
	"log/slog"
)

type GenerateTicketRequest struct{}
type GenerateTicketResponse struct{}

type Handler struct {
	logger  *slog.Logger
	wpc     *weasyPrintClient.WeasyPrintClient
	outPath string
}

func New(logger *slog.Logger, wpc *weasyPrintClient.WeasyPrintClient, outPath string) *Handler {
	return &Handler{
		logger,
		wpc,
		outPath,
	}
}

func (h *Handler) GenerateTicketHandler(_ context.Context, input *GenerateTicketRequest) (*GenerateTicketResponse, error) {
	doc, err := h.createTicket(input)
	if err != nil {
		return nil, err
	}
	err = utils.WriteToFile(doc, h.outPath, "test-ticket.pdf")
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	return nil, nil
}

func (h *Handler) createTicket(_ *GenerateTicketRequest) (*bytes.Buffer, error) {
	doc, err := utils.RenderTemplate(templates.Ticket())
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	var pdf bytes.Buffer
	err = h.wpc.GeneratePDF(doc, &pdf)
	if err != nil {
		h.logger.Error(err.Error())
		return nil, err
	}
	return &pdf, nil

}
