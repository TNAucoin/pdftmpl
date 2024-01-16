package router

import (
	"context"
	"fmt"
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
		OperationID: "get-health",
		Summary:     "Health Check",
		Method:      http.MethodGet,
		Path:        "/livez",
	}, func(ctx context.Context, input *health.HealthInput) (*health.HealthOutput, error) {
		resp := health.GetHealth()
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "generator-create",
		Summary:     "Generate Document Data",
		Method:      http.MethodPost,
		Path:        "/v1/generate",
	}, func(ctx context.Context, input *generator.GenerateInput) (*generator.GenerateOutput, error) {
		fmt.Println(ctx.Value("Method"))
		resp, err := generator.Create(logger, gotenClient)(input)
		if err != nil {
			return nil, err
		}
		return resp, nil
	})

	return r
}
