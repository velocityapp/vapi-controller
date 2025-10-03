package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type EchoHandler struct{}

func (echo EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestHeaders := make(map[string]interface{})

	for name, values := range r.Header {
		requestHeaders[name] = values
	}

	path := r.URL.Path
	method := r.Method

	response := make(map[string]interface{})

	response["url"] = path
	response["method"] = method
	response["request-headers"] = requestHeaders
	response["datetime"] = time.Now().String()

	// Echo body
	//body, err := io.ReadAll(r.Body)

	body := make(map[string]interface{})

	bodyErr := json.NewDecoder(r.Body).Decode(&body)
	if bodyErr != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	response["body"] = body

	slog.Debug("Interface", "response", response)
	//encodeErr := json.NewEncoder(w).Encode(response)
	//buf := new(bytes.Buffer)

	contentType := "application/json"
	w.Header().Add("content-type", contentType)
	//len, err := w.Write(buf)

	//buf, err := json.Marshal(response)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("error encoding json response", "err", err)
	}
}
