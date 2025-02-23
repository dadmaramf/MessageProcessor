// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	context "context"
	model "messageprocessor/internal/model"
	reflect "reflect"
	time "time"

	sarama "github.com/IBM/sarama"
	gomock "github.com/golang/mock/gomock"
)

// MockSyncProducer is a mock of SyncProducer interface.
type MockSyncProducer struct {
	ctrl     *gomock.Controller
	recorder *MockSyncProducerMockRecorder
}

// MockSyncProducerMockRecorder is the mock recorder for MockSyncProducer.
type MockSyncProducerMockRecorder struct {
	mock *MockSyncProducer
}

// NewMockSyncProducer creates a new mock instance.
func NewMockSyncProducer(ctrl *gomock.Controller) *MockSyncProducer {
	mock := &MockSyncProducer{ctrl: ctrl}
	mock.recorder = &MockSyncProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSyncProducer) EXPECT() *MockSyncProducerMockRecorder {
	return m.recorder
}

// AbortTxn mocks base method.
func (m *MockSyncProducer) AbortTxn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AbortTxn")
	ret0, _ := ret[0].(error)
	return ret0
}

// AbortTxn indicates an expected call of AbortTxn.
func (mr *MockSyncProducerMockRecorder) AbortTxn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AbortTxn", reflect.TypeOf((*MockSyncProducer)(nil).AbortTxn))
}

// AddMessageToTxn mocks base method.
func (m *MockSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessageToTxn", msg, groupId, metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMessageToTxn indicates an expected call of AddMessageToTxn.
func (mr *MockSyncProducerMockRecorder) AddMessageToTxn(msg, groupId, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessageToTxn", reflect.TypeOf((*MockSyncProducer)(nil).AddMessageToTxn), msg, groupId, metadata)
}

// AddOffsetsToTxn mocks base method.
func (m *MockSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOffsetsToTxn", offsets, groupId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOffsetsToTxn indicates an expected call of AddOffsetsToTxn.
func (mr *MockSyncProducerMockRecorder) AddOffsetsToTxn(offsets, groupId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOffsetsToTxn", reflect.TypeOf((*MockSyncProducer)(nil).AddOffsetsToTxn), offsets, groupId)
}

// BeginTxn mocks base method.
func (m *MockSyncProducer) BeginTxn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTxn")
	ret0, _ := ret[0].(error)
	return ret0
}

// BeginTxn indicates an expected call of BeginTxn.
func (mr *MockSyncProducerMockRecorder) BeginTxn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTxn", reflect.TypeOf((*MockSyncProducer)(nil).BeginTxn))
}

// Close mocks base method.
func (m *MockSyncProducer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockSyncProducerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSyncProducer)(nil).Close))
}

// CommitTxn mocks base method.
func (m *MockSyncProducer) CommitTxn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitTxn")
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitTxn indicates an expected call of CommitTxn.
func (mr *MockSyncProducerMockRecorder) CommitTxn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitTxn", reflect.TypeOf((*MockSyncProducer)(nil).CommitTxn))
}

// IsTransactional mocks base method.
func (m *MockSyncProducer) IsTransactional() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTransactional")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTransactional indicates an expected call of IsTransactional.
func (mr *MockSyncProducerMockRecorder) IsTransactional() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTransactional", reflect.TypeOf((*MockSyncProducer)(nil).IsTransactional))
}

// SendMessage mocks base method.
func (m *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", msg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockSyncProducerMockRecorder) SendMessage(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockSyncProducer)(nil).SendMessage), msg)
}

// SendMessages mocks base method.
func (m *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessages", msgs)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessages indicates an expected call of SendMessages.
func (mr *MockSyncProducerMockRecorder) SendMessages(msgs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessages", reflect.TypeOf((*MockSyncProducer)(nil).SendMessages), msgs)
}

// TxnStatus mocks base method.
func (m *MockSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxnStatus")
	ret0, _ := ret[0].(sarama.ProducerTxnStatusFlag)
	return ret0
}

// TxnStatus indicates an expected call of TxnStatus.
func (mr *MockSyncProducerMockRecorder) TxnStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxnStatus", reflect.TypeOf((*MockSyncProducer)(nil).TxnStatus))
}

// MockConsumerGroup is a mock of ConsumerGroup interface.
type MockConsumerGroup struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerGroupMockRecorder
}

// MockConsumerGroupMockRecorder is the mock recorder for MockConsumerGroup.
type MockConsumerGroupMockRecorder struct {
	mock *MockConsumerGroup
}

// NewMockConsumerGroup creates a new mock instance.
func NewMockConsumerGroup(ctrl *gomock.Controller) *MockConsumerGroup {
	mock := &MockConsumerGroup{ctrl: ctrl}
	mock.recorder = &MockConsumerGroupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerGroup) EXPECT() *MockConsumerGroupMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConsumerGroup) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConsumerGroupMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConsumerGroup)(nil).Close))
}

// Consume mocks base method.
func (m *MockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", ctx, topics, handler)
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume.
func (mr *MockConsumerGroupMockRecorder) Consume(ctx, topics, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockConsumerGroup)(nil).Consume), ctx, topics, handler)
}

// Errors mocks base method.
func (m *MockConsumerGroup) Errors() <-chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Errors")
	ret0, _ := ret[0].(<-chan error)
	return ret0
}

// Errors indicates an expected call of Errors.
func (mr *MockConsumerGroupMockRecorder) Errors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errors", reflect.TypeOf((*MockConsumerGroup)(nil).Errors))
}

// Pause mocks base method.
func (m *MockConsumerGroup) Pause(partitions map[string][]int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Pause", partitions)
}

// Pause indicates an expected call of Pause.
func (mr *MockConsumerGroupMockRecorder) Pause(partitions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pause", reflect.TypeOf((*MockConsumerGroup)(nil).Pause), partitions)
}

// PauseAll mocks base method.
func (m *MockConsumerGroup) PauseAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PauseAll")
}

// PauseAll indicates an expected call of PauseAll.
func (mr *MockConsumerGroupMockRecorder) PauseAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseAll", reflect.TypeOf((*MockConsumerGroup)(nil).PauseAll))
}

// Resume mocks base method.
func (m *MockConsumerGroup) Resume(partitions map[string][]int32) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Resume", partitions)
}

// Resume indicates an expected call of Resume.
func (mr *MockConsumerGroupMockRecorder) Resume(partitions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resume", reflect.TypeOf((*MockConsumerGroup)(nil).Resume), partitions)
}

// ResumeAll mocks base method.
func (m *MockConsumerGroup) ResumeAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResumeAll")
}

// ResumeAll indicates an expected call of ResumeAll.
func (mr *MockConsumerGroupMockRecorder) ResumeAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResumeAll", reflect.TypeOf((*MockConsumerGroup)(nil).ResumeAll))
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// SaveMessage mocks base method.
func (m *MockService) SaveMessage(msg string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMessage", msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMessage indicates an expected call of SaveMessage.
func (mr *MockServiceMockRecorder) SaveMessage(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMessage", reflect.TypeOf((*MockService)(nil).SaveMessage), msg)
}

// SentMessages mocks base method.
func (m *MockService) SentMessages() ([]model.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SentMessages")
	ret0, _ := ret[0].([]model.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SentMessages indicates an expected call of SentMessages.
func (mr *MockServiceMockRecorder) SentMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SentMessages", reflect.TypeOf((*MockService)(nil).SentMessages))
}

// StartConsumerProcessingMessage mocks base method.
func (m *MockService) StartConsumerProcessingMessage(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartConsumerProcessingMessage", ctx)
}

// StartConsumerProcessingMessage indicates an expected call of StartConsumerProcessingMessage.
func (mr *MockServiceMockRecorder) StartConsumerProcessingMessage(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartConsumerProcessingMessage", reflect.TypeOf((*MockService)(nil).StartConsumerProcessingMessage), ctx)
}

// StartProcessingMessage mocks base method.
func (m *MockService) StartProcessingMessage(ctx context.Context, handlePeriod time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartProcessingMessage", ctx, handlePeriod)
}

// StartProcessingMessage indicates an expected call of StartProcessingMessage.
func (mr *MockServiceMockRecorder) StartProcessingMessage(ctx, handlePeriod interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartProcessingMessage", reflect.TypeOf((*MockService)(nil).StartProcessingMessage), ctx, handlePeriod)
}
