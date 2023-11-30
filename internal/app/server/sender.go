package server

import (
	"encoding/json"
	"fmt"
	"homework6/internal/infrastructure/kafka"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
)

type ServerMessage struct {
	CreationTime  time.Time
	MethodType    string
	RequestMethod string
	RequestBody   string
}

type KafkaSender struct {
	producer kafka.ProducerKafka
	topic    string
}

func NewKafkaSender(producer kafka.ProducerKafka, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func NewServerMessage(req *http.Request, reqBody []byte) (*ServerMessage, error) {
	return &ServerMessage{
		CreationTime:  time.Now(),
		MethodType:    req.URL.Path,
		RequestMethod: req.Method,
		RequestBody:   string(reqBody),
	}, nil

}

func (s *KafkaSender) SendMessage(message *sarama.ProducerMessage) error {

	partition, offset, err := s.producer.SendSyncMessage(message)

	if err != nil {
		fmt.Println("Send message connector error", err)
		return err
	}

	fmt.Println("Send message:\nPartition: ", partition, " Offset: ", offset, " Type:", message.Key)
	return nil
}

func BuildMessage(message *ServerMessage, topic string) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		fmt.Println("Send message marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(message.MethodType)),
	}, nil
}
