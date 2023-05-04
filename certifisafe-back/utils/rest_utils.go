package utils

import (
	"encoding/json"
	"net/http"
)

func ReturnResponse(w http.ResponseWriter, err error, data interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "error when encoding json", http.StatusInternalServerError)
		return
	}
}
