package services_test

import (
	"context"
	"errors"
	"log/slog"
	"messageprocessor/internal/model"
	messagereader "messageprocessor/internal/services/message_reader"
	messagesender "messageprocessor/internal/services/message_sender"
	mock_services "messageprocessor/internal/services/mock"
	mock_storage "messageprocessor/internal/storage/mock"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
)

func TestStartProcessingMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock_storage.NewMockStorage(ctrl)
	mockProducer := mock_services.NewMockSyncProducer(ctrl)

	logger := slog.Default()

	sender := messagesender.New(mockStorage, mockProducer, logger)

	handlePeriod := 10 * time.Millisecond

	tests := []struct {
		name          string
		message       *model.Message
		getNewOutbox  error
		encodeError   error
		sendError     error
		setDownError  error
		expectedError error
	}{
		{
			name:          "Success",
			message:       &model.Message{ID: 1, Content: "test message"},
			getNewOutbox:  nil,
			encodeError:   nil,
			sendError:     nil,
			setDownError:  nil,
			expectedError: nil,
		},
		{
			name:          "GetNewOutboxError",
			message:       nil,
			getNewOutbox:  errors.New("get new outbox error"),
			encodeError:   nil,
			sendError:     nil,
			setDownError:  nil,
			expectedError: errors.New("get new outbox error"),
		},
		{
			name:          "EncodeError",
			message:       &model.Message{ID: 1, Content: "test message"},
			getNewOutbox:  nil,
			encodeError:   errors.New("encode error"),
			sendError:     nil,
			setDownError:  nil,
			expectedError: errors.New("encode error"),
		},
		{
			name:          "SendError",
			message:       &model.Message{ID: 1, Content: "test message"},
			getNewOutbox:  nil,
			encodeError:   nil,
			sendError:     errors.New("send error"),
			setDownError:  nil,
			expectedError: errors.New("send error"),
		},
		{
			name:          "SetDownError",
			message:       &model.Message{ID: 1, Content: "test message"},
			getNewOutbox:  nil,
			encodeError:   nil,
			sendError:     nil,
			setDownError:  errors.New("set down error"),
			expectedError: errors.New("set down error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mockStorage.EXPECT().GetNewOutbox(gomock.Any()).Return(test.message, test.getNewOutbox).AnyTimes()

			if test.message != nil {
				mockProducer.EXPECT().SendMessage(gomock.Any()).Return(int32(0), int64(0), test.sendError).AnyTimes()

				if test.sendError == nil {
					mockStorage.EXPECT().SetDown(test.message.ID).Return(test.setDownError).AnyTimes()
				}
			}

			go sender.StartProcessingMessage(ctx, handlePeriod)

			time.Sleep(2 * handlePeriod)

			cancel()
		})
	}
}

type MockConsumerGroup struct {
	mock.Mock
}

func (m *MockConsumerGroup) PauseAll() {
	panic("unimplemented")
}

func (m *MockConsumerGroup) ResumeAll() {
	panic("unimplemented")
}

func (m *MockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)
	return args.Error(0)
}

func (m *MockConsumerGroup) Errors() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func (m *MockConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}
func TestStartConsumerProcessingMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockConsumer := new(MockConsumerGroup)
	mockLogger := slog.Default()
	mockStorage := mock_storage.NewMockStorage(ctrl)

	reader := messagereader.New(mockConsumer, mockStorage, mockLogger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockConsumer.On("Consume", mock.Anything, []string{"processedMessage"}, mock.Anything).Return(nil)

	reader.StartConsumerProcessingMessage(ctx)

	time.Sleep(100 * time.Millisecond)

	mockConsumer.AssertCalled(t, "Consume", mock.Anything, []string{"processedMessage"}, mock.Anything)
}
