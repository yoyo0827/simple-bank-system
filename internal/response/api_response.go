package response

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// 成功
func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ApiResponse{
		Status: "success",
		Data:   data,
	})
}

// 錯誤
func WriteError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ApiResponse{
		Status: "error",
		Error:  errMsg,
	})
}
