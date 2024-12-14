package response

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondWithJSON(w http.ResponseWriter, code int, status bool, message string, payload interface{}) {
	responseStruct := JSONResponse{
		Status:  status,
		Message: message,
		Data:    payload,
	}

	response, err := json.Marshal(responseStruct)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
