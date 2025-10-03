package routes

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

type compressResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (crw *compressResponseWriter) Write(data []byte) (int, error) {
	return crw.writer.Write(data)
}

// ResponseCompressor will compress response if client supports it
func ResponseCompressor(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gz := gzip.NewWriter(w)
		w.Header().Set("Content-Encoding", "gzip")
		next(&compressResponseWriter{w, gz}, r)
		gz.Close()
	} else {
		next(w, r)
	}
}

func Recover(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	defer func() {
		if rec := recover(); rec != nil {
			slog.Error("Panic recovered: %v\nStack trace:\n%s", rec, debug.Stack())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"}); err != nil {
				slog.Error("Error ecoding response", "err", err)
			}
		}
	}()
	next(w, r)
}
