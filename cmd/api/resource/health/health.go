package health

import (
	"context"
)

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	Body struct {
		Message string `json:"message" example:"livez" doc:"Status message"`
	}
}

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheckHandler(_ context.Context, input *HealthCheckRequest) (*HealthCheckResponse, error) {
	resp := h.getHealth()
	return resp, nil
}

func (h *Handler) getHealth() *HealthCheckResponse {
	resp := &HealthCheckResponse{}
	resp.Body.Message = "livez"
	return resp
}
