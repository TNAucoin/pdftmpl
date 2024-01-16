package health

func GetHealth() *HealthOutput {
	resp := &HealthOutput{}
	resp.Body.Message = "livez"
	return resp
}
