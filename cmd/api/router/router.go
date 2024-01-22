package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/generate"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/health"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/ticket"
	"github.com/tnaucoin/pdftmpl/config"
	"github.com/tnaucoin/pdftmpl/internal/weasyPrintClient"
	"log/slog"
	"net/http"
)

func New(logger *slog.Logger, wpc weasyPrintClient.WeasyPrintClient, cfg *config.Conf) *chi.Mux {
	healthHandler := health.New()
	invoiceGeneratorHandler := generate.New(logger, wpc, cfg.Server.VolumeOutPath)
	ticketGeneratorHandler := ticket.New(logger, wpc, cfg.Server.VolumeOutPath)
	// Create a new router and API instance.
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("PDF Generator API", "1.0.0"))
	// Register the health check and invoice generator handlers.
	registerHealthCheck(api, healthHandler)
	registerInvoiceGenerator(api, invoiceGeneratorHandler)
	registerTicketGenerator(api, ticketGeneratorHandler)

	return r
}

// registerHealthCheck registers a health check handler in the API.
func registerHealthCheck(api huma.API, handler *health.Handler) {
	huma.Register(api, huma.Operation{
		OperationID:   "get-health",
		Summary:       "Health Check",
		Method:        http.MethodGet,
		Path:          "/livez",
		DefaultStatus: http.StatusOK,
	}, handler.HealthCheckHandler)
}

// registerInvoiceGenerator registers an invoice generator handler in the API.
func registerInvoiceGenerator(api huma.API, handler *generate.Handler) {
	huma.Register(api, huma.Operation{
		OperationID:   "generator-create",
		Summary:       "Generate Document Data",
		Method:        http.MethodPost,
		Path:          "/v1/generate/invoice",
		DefaultStatus: http.StatusCreated,
	}, handler.GenerateInvoicePDFHandler)
}

func registerTicketGenerator(api huma.API, handler *ticket.Handler) {
	huma.Register(api, huma.Operation{
		OperationID:   "ticket-create",
		Summary:       "Generate Ticket",
		Method:        http.MethodPost,
		Path:          "/v1/generate/ticket",
		DefaultStatus: http.StatusCreated,
	}, handler.GenerateTicketHandler)
}
