package server

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewServerMessage(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		jsonBody := []byte(`{"name":"asd","rating":10}`)
		req := httptest.NewRequest("POST", "http://localhost:9000/article", nil)
		expectedServerMessage := &ServerMessage{
			MethodType:    "/article",
			RequestMethod: "POST",
			RequestBody:   string(`{"name":"asd","rating":10}`),
		}

		getMessage, err := NewServerMessage(req, jsonBody)
		getMessage.CreationTime = time.Time{}

		require.Nil(t, err)
		assert.Equal(t, expectedServerMessage, getMessage)

	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest("GET", "http://localhost:9000/article?id=4", nil)
		expectedServerMessage := &ServerMessage{
			MethodType:    "/article",
			RequestMethod: "GET",
		}

		getMessage, err := NewServerMessage(req, nil)
		getMessage.CreationTime = time.Time{}

		require.Nil(t, err)
		assert.Equal(t, expectedServerMessage, getMessage)

	})
}

func Test_BuildMessage(t *testing.T) {
	t.Parallel()
	topic := "requests"
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sendMessage := &ServerMessage{
			MethodType:    "/article",
			RequestMethod: "POST",
			RequestBody:   string(`{"name":"asd","rating":10}`),
		}
		expectedValue := []byte(`{"CreationTime":"0001-01-01T00:00:00Z","MethodType":"/article","RequestMethod":"POST","RequestBody":"{\"name\":\"asd\",\"rating\":10}"}`)
		expectedMessage := &sarama.ProducerMessage{
			Topic:     topic,
			Value:     sarama.ByteEncoder(expectedValue),
			Partition: -1,
			Key:       sarama.StringEncoder("/article"),
		}

		gotMessage, err := BuildMessage(sendMessage, topic)

		require.Nil(t, err)
		require.Equal(t, expectedMessage.Topic, gotMessage.Topic)
		require.Equal(t, expectedMessage.Key, gotMessage.Key)
		value, err := gotMessage.Value.Encode()
		require.Nil(t, err)
		require.Equal(t, string(expectedValue), string(value))

	})

}

func Test_SendMessage(t *testing.T) {
	topic := "requests"
	value := []byte(`{"CreationTime":"0001-01-01T00:00:00Z","MethodType":"/article","RequestMethod":"POST","RequestBody":"{\"name\":\"asd\",\"rating\":10}"}`)
	sendMessage := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(value),
		Partition: -1,
		Key:       sarama.StringEncoder("/article"),
	}
	t.Run("success", func(t *testing.T) {

		sender := setUp(t, topic)
		defer sender.tearDown()
		sender.Producer.EXPECT().SendSyncMessage(sendMessage).Return(int32(0), int64(0), nil)

		err := sender.Sender.SendMessage(sendMessage)
		require.Nil(t, err)

	})

	t.Run("success", func(t *testing.T) {

		sender := setUp(t, topic)
		defer sender.tearDown()

		sender.Producer.EXPECT().SendSyncMessage(sendMessage).Return(int32(0), int64(0), errors.New("some error"))

		err := sender.Sender.SendMessage(sendMessage)
		require.Error(t, err)

	})

}
