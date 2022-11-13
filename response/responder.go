package response

import (
	"encoding/json"
	"net/http"
)

func WriteJson(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		panic(err)
	}
}
func WriteXml(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Add("Content-Type", "application/xml")
	writer.WriteHeader(code)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		panic(err)
	}
}
