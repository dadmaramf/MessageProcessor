package app

import (
	"log/slog"
	"messageprocessor/internal/config"
	"messageprocessor/internal/services"
	"net/http"
)

func NewServer(logger *slog.Logger, cfg *config.Config, service services.Service) http.Handler {
	mux := http.NewServeMux()
	addRouters(mux, logger, service)
	var handler http.Handler = mux
	return handler

}
