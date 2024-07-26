package api

import (
	"log/slog"
	"messageprocessor/internal/app/handlers"
	"messageprocessor/internal/config"
	"messageprocessor/internal/services"
	"net/http"
)

func addRouters(mux *http.ServeMux, log *slog.Logger, cfg *config.Config, service services.Service) {
	mux.HandleFunc("/message", handlers.SaveMessage(log, service))
}
