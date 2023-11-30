// Code generated by MockGen. DO NOT EDIT.
// Source: ./kafka.go

// Package mock_kafka is a generated GoMock package.
package mock_kafka

import (
	reflect "reflect"

	sarama "github.com/Shopify/sarama"
	gomock "github.com/golang/mock/gomock"
)

// MockProducerKafka is a mock of ProducerKafka interface.
type MockProducerKafka struct {
	ctrl     *gomock.Controller
	recorder *MockProducerKafkaMockRecorder
}

// MockProducerKafkaMockRecorder is the mock recorder for MockProducerKafka.
type MockProducerKafkaMockRecorder struct {
	mock *MockProducerKafka
}

// NewMockProducerKafka creates a new mock instance.
func NewMockProducerKafka(ctrl *gomock.Controller) *MockProducerKafka {
	mock := &MockProducerKafka{ctrl: ctrl}
	mock.recorder = &MockProducerKafkaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducerKafka) EXPECT() *MockProducerKafkaMockRecorder {
	return m.recorder
}

// SendSyncMessage mocks base method.
func (m *MockProducerKafka) SendSyncMessage(message *sarama.ProducerMessage) (int32, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSyncMessage", message)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SendSyncMessage indicates an expected call of SendSyncMessage.
func (mr *MockProducerKafkaMockRecorder) SendSyncMessage(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSyncMessage", reflect.TypeOf((*MockProducerKafka)(nil).SendSyncMessage), message)
}