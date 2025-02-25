package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"transactionService/database"
	"transactionService/kafka"

	"github.com/IBM/sarama"
	"github.com/hashicorp/go-hclog"
)

func main() {

	fmt.Println("Starting to develop banking application")
	// logging app file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Some error occured in creating or opeing a log file", err)
		os.Exit(1)
	}

	// logger
	loggs := hclog.New(&hclog.LoggerOptions{
		Name:       "Banking App",
		Output:     logFile,
		Level:      hclog.Debug,
		JSONFormat: false,
	})

	loggs.Info("Hell World")

	dsn := "postgres://postgres:abcd@postgres:5432/accounts?sslmode=disable"

	db := database.NewPostgresPoolDB(dsn, 10, 2)
	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	defer db.Close(ctx)

	// handlers
	kafkaConsumer := kafka.NewKafkaConsumer(db)

	// Kafka configuration
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	brokers := []string{"kafka:9092"}
	groupID := "transaction-group"
	topics := []string{"transaction"}

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to start consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// Reconnect and resume on errors
			err = consumerGroup.Consume(ctx, topics, kafkaConsumer)
			if err != nil {
				log.Printf("Consumer error: %v", err)
			}
			if ctx.Err() != nil {
				return // Exit if context is canceled
			}
		}
	}()

	// Listen for errors
	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("Consumer group error: %v", err)
		}
	}()

	log.Println("Consumer group started. Waiting for messages...")

	// Handle SIGINT/SIGTERM for shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down consumer...")
	cancel()
	wg.Wait()
}
