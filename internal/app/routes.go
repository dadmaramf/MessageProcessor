package app

import (
	"log/slog"
	"messageprocessor/internal/app/handlers"
	"messageprocessor/internal/services"
	"net/http"
)

func addRouters(mux *http.ServeMux, log *slog.Logger, service services.Service) {
	mux.HandleFunc("/submit", handlers.SaveMessage(log, service))
	mux.HandleFunc("/state", handlers.SentMessages(log, service))
}
