package types

type HealthInput struct{}

type HealthOutput struct {
	Body struct {
		Message string `json:"message" example:"livez" doc:"Status message"`
	}
}
