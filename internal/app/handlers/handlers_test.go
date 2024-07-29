package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"messageprocessor/internal/app/handlers"
	"messageprocessor/internal/model"
	mock_services "messageprocessor/internal/services/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

func TestSaveMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockService(ctrl)
	logger := slog.Default()

	handler := handlers.SaveMessage(logger, mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		saveMessageErr error
	}{
		{
			name: "Success",
			requestBody: map[string]string{
				"message": "test message",
			},
			expectedStatus: http.StatusCreated,
			saveMessageErr: nil,
		},
		{
			name:           "InvalidRequest",
			requestBody:    "invalid request body",
			expectedStatus: http.StatusBadRequest,
			saveMessageErr: nil,
		},
		{
			name: "SaveMessageError",
			requestBody: map[string]string{
				"message": "test message",
			},
			expectedStatus: http.StatusInternalServerError,
			saveMessageErr: errors.New("failed to save message"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/message", bytes.NewReader(body))
			w := httptest.NewRecorder()

			if test.expectedStatus != http.StatusBadRequest {
				mockService.EXPECT().SaveMessage(gomock.Any()).Return(test.saveMessageErr)
			}

			handler(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}

func TestSentMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_services.NewMockService(ctrl)
	logger := slog.Default()

	handler := handlers.SentMessages(logger, mockService)

	tests := []struct {
		name           string
		expectedStatus int
		messages       []model.Message
		getMessagesErr error
	}{
		{
			name:           "Success",
			expectedStatus: http.StatusOK,
			messages:       []model.Message{{ID: 1, Content: "test message"}},
			getMessagesErr: nil,
		},
		{
			name:           "GetMessagesError",
			expectedStatus: http.StatusInternalServerError,
			messages:       nil,
			getMessagesErr: errors.New("failed to get messages"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/message/state", nil)
			w := httptest.NewRecorder()

			mockService.EXPECT().SentMessages().Return(test.messages, test.getMessagesErr)

			handler(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
