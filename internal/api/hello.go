package api

import (
	"encoding/json"
	"net/http"
)

// HelloHandler godoc
// @Summary Hello
// @Description Returns a hello message
// @Tags hello
// @Produce json
// @Success 200 {object} map[string]string
// @Router /hello [get]
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, Bank System 👋"})
}
