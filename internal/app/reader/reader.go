package reader

import (
	"context"
	"homework6/internal/infrastructure/kafka"
	"log"
	"sync"

	"github.com/Shopify/sarama"
)

func StartReader(ctx context.Context, clientConsumer sarama.ConsumerGroup) {

	consumer := kafka.NewConsumerGroup()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := clientConsumer.Consume(ctx, []string{"requests"}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	log.Println("Sarama consumer up and running!...")

	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	}
	
	wg.Wait()
}
