// Package helpers implements commonly used functions (response API)//
package helpers

import (
	"encoding/json"
	"net/http"
)

// APIResponse is
type APIResponse struct {
	Message interface{} `json:"message"`
}

// JSONResponse represent a json response
type JSONResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

// Response is
func Response(w http.ResponseWriter, httpStatus int, data interface{}) {
	apiResponse := new(APIResponse)
	apiResponse.Message = data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiResponse)
}

// DirectResponse is
func DirectResponse(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}
