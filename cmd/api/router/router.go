package router

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/generator"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/health"
	"github.com/tnaucoin/pdftmpl/internal/pdfClient"
	"log/slog"
	"net/http"
)

func New(logger *slog.Logger, gotenClient *pdfClient.PdfClient) *chi.Mux {
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("PDF Generator API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID:   "get-health",
		Summary:       "Health Check",
		Method:        http.MethodGet,
		Path:          "/livez",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *health.HealthInput) (*health.HealthOutput, error) {
		resp := health.GetHealth()
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "generator-create",
		Summary:       "Generate Document Data",
		Method:        http.MethodPost,
		Path:          "/v1/generate",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *generator.GenerateInput) (*generator.GenerateOutput, error) {
		err := generator.Create(logger, gotenClient)(input)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		return nil, nil
	})

	return r
}
