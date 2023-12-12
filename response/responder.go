package response

import (
	"encoding/json"
	"net/http"
)

// WriteJson writes the data to the ResponseWriter with Content-Type as application/json and the code as a header.
func WriteJson(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		panic(err)
	}
}

// WriteXml writes the data to the ResponseWriter with Content-Type as application/xml and the code as a header.
func WriteXml(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Add("Content-Type", "application/xml")
	writer.WriteHeader(code)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		panic(err)
	}
}
