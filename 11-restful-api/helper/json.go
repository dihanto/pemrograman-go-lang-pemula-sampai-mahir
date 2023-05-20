package helper

import (
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	err := json.NewDecoder(request.Body).Decode(result)
	PanifIfError(err)
}

func SendToWriter(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(response)
	PanifIfError(err)
}
