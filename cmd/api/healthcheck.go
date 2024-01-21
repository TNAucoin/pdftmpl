package main

import "github.com/tnaucoin/pdftmpl/cmd/api/resource/types"

func GetHealth() *types.HealthOutput {
	resp := &types.HealthOutput{}
	resp.Body.Message = "livez"
	return resp
}
