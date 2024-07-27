package handlers

import (
	"log/slog"
	"messageprocessor/internal/services"
	"net/http"
)

// POST/message
func SaveMessage(log *slog.Logger, service services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.handler.SaveMessage"
		log.With(slog.String("op", op))
		resp := struct {
			Message string `json:"message"`
		}{}

		if err := decode(r, &resp); err != nil {
			log.Error("Failed to decode request", "error", err)
			errorJSON(w, http.StatusBadRequest, "Failed to decode request")
			return
		}

		if err := service.SaveMessage(resp.Message); err != nil {
			log.Error("Failed to save message", "error", err)
			errorJSON(w, http.StatusInternalServerError, "Failed to save message")
			return
		}

		successMsg := &struct {
			Answer string `json:"answer"`
		}{"Message received successfully"}

		if err := encode(w, http.StatusCreated, successMsg); err != nil {
			log.Error("Failed to create answer", "error", err)
			errorJSON(w, http.StatusInternalServerError, "Failed to create answer")
			return
		}
		log.Info("request body decoded", slog.Any("request", successMsg))
	}
}

// GET/message/state
func SentMessages(log *slog.Logger, service services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.handler.SentMessages"
		log.With(slog.String("op", op))

		messages, err := service.SentMessages()

		if err != nil {
			log.Error("Failed to get messages", "error", err)
			errorJSON(w, http.StatusInternalServerError, "Failed to get messages")
			return
		}

		if err = encode(w, http.StatusOK, messages); err != nil {
			log.Error("Failed to create answer", "error", err)
			errorJSON(w, http.StatusInternalServerError, "Failed to create answer")
			return
		}
		log.Info("request body decoded", slog.Any("request", messages))
	}
}
