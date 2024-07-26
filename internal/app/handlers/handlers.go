package handlers

import (
	"log/slog"
	"messageprocessor/internal/services"
	"net/http"
)

func SaveMessage(log *slog.Logger, service services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.handler.SaveMessage"
		var message string
		if err := decode(r, &message); err != nil {
			log.Error("Failed to decode request", "op", op, "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
		}

		if err := service.SaveMessage(message); err != nil {
			log.Error("Failed to save message", "op", op, "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		successMsg := &struct{ answer string }{"Message received successfully"}
		if err := encode(w, http.StatusCreated, successMsg); err != nil {
			log.Error("Failed to create answer", "op", op, "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func SentMessages(log *slog.Logger, service services.Service) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			const op = "api.handler.SentMessages"

			messages, err := service.SentMessages()

			if err != nil {
				log.Error("Failed to get messages", "op", op, "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if err = encode(w, http.StatusOK, messages); err != nil {
				log.Error("Failed to create answer", "op", op, "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
