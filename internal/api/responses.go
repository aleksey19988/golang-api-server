package api

import (
	"encoding/json"
	"net/http"
)

func writeErrorResponse(
	w http.ResponseWriter,
	data string,
	code int,
) {
	w.WriteHeader(code)
	response := Response{
		Status: StatusError,
		Code:   code,
		Data:   data,
	}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}

func writeSuccessResponse(w http.ResponseWriter, data string, code int) {
	w.WriteHeader(code)
	response := Response{
		Status: StatusSuccess,
		Code:   code,
		Data:   data,
	}
	responseJSON, _ := json.Marshal(response)
	w.Write(responseJSON)
}
