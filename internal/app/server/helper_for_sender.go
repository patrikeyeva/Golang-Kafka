package server

import (
	mock_kafka "homework6/internal/infrastructure/kafka/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

type senderFixture struct {
	ctrl     *gomock.Controller
	Producer *mock_kafka.MockProducerKafka
	Sender   *KafkaSender
}

func setUp(t *testing.T, topic string) senderFixture {
	ctrl := gomock.NewController(t)
	producer := mock_kafka.NewMockProducerKafka(ctrl)
	s := NewKafkaSender(producer, topic)

	return senderFixture{
		ctrl:     ctrl,
		Producer: producer,
		Sender:   s,
	}
}

func (s *senderFixture) tearDown() {
	s.ctrl.Finish()
}
