package utils

import (
	"encoding/json"
	"math"
	"net/http"
)

// JSONResponse will accept a payload object and status code
// the function will set code as the HTTP status code for the response writer
// the function will set the payload as the body for the response write
func JSONResponse(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ToFixed will take a float value and a precision value to set the float to the specified precision
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
