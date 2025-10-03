package routes

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"vapp/controller/config"
)

var GlobalStr = ""

func init() {
	count := 1024 * 8
	for i := 0; i < count; i++ {
		GlobalStr += "1"
	}
}

func Start() {

	serverPort, err := config.GetConfig[string]("server.port")
	if err != nil {
		slog.Error("Failed to load server port from config")
	}
	serverAddress, err := config.GetConfig[string]("server.address")
	if err != nil {
		slog.Error("Failed to load server address from config")
	}

	addr := *serverAddress + ":" + *serverPort

	ctx := context.Background()

	s := &http.Server{
		Addr:           addr,
		BaseContext:    func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 4,
	}

	//mux := http.NewServeMux()
	mux := new(VMux)
	mux.Add(ResponseCompressor)
	mux.Add(Recover) // has to be last...

	echoHandler := new(EchoHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.Handle("/", echoHandler)
	s.Handler = mux

	slog.Info("Starting server", "address", addr)
	srvErr := make(chan error, 1)

	go func() {
		srvErr <- s.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		slog.Error("Server error", "error", err)
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		slog.Info("Shutting down server")
		s.Shutdown(ctx)
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = s.Shutdown(context.Background())
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server shutdown error", "error", err)
		return
	}

	slog.Info("Server stopped")
}
