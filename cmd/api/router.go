package main

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/tnaucoin/pdftmpl/cmd/api/resource/types"
	"net/http"
)

func (app *Application) Router() *chi.Mux {
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("PDF Generator API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID:   "get-health",
		Summary:       "Health Check",
		Method:        http.MethodGet,
		Path:          "/livez",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, input *types.HealthInput) (*types.HealthOutput, error) {
		resp := GetHealth()
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "generator-create",
		Summary:       "Generate Document Data",
		Method:        http.MethodPost,
		Path:          "/v1/generate",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *types.GenerateInput) (*types.GenerateOutput, error) {
		err := generatePdfPost(app)(input)
		if err != nil {
			app.logger.Error(err.Error())
			return nil, err
		}
		return nil, nil
	})

	return r
}
