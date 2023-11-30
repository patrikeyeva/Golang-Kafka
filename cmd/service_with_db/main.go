package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"homework6/internal/app/reader"
	"homework6/internal/app/server"
	"homework6/internal/infrastructure/kafka"
	"homework6/internal/pkg/db"
	"homework6/internal/pkg/repository/postgresql"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	brokersStr := os.Getenv("KAFKA_BROKERS")
	brokers := strings.Split(brokersStr, ",")

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientConsumer, err := sarama.NewConsumerGroup(brokers, "service-with-db", kafka.CreateConfigForConsumerGroup())
	if err != nil {
		log.Fatal("Error creating consumer group client:", err)
	}
	defer func() {
		if err := clientConsumer.Close(); err != nil {
			log.Panicf("Error closing client: %v", err)
		}
	}()

	database, err := db.NewDBWithDSN(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	mainServer := server.Server{
		ArticleRepo: postgresql.NewArticles(database),
		CommentRepo: postgresql.NewComments(database),
		Sender:      server.NewKafkaSender(kafkaProducer, "requests"),
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		http.Handle("/", server.CreateRouter(ctx, mainServer))
		if err := http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		reader.StartReader(ctx, clientConsumer)
	}()

	wg.Wait()

}
