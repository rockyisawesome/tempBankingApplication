package main

import (
	"context"
	"fmt"
	"ledgerservice/configurations"
	"ledgerservice/database"
	"ledgerservice/kafka"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	// getting mongodb configurations
	mongodbconfig, err := configurations.NewMongoDbConfig()
	if err != nil {
		loggs.Error("Not able to create Retrieve App Configurations")
		os.Exit(1)
	}

	mongodb := database.NewMongoDB(mongodbconfig, &loggs)
	ctx := context.Background()

	if err := mongodb.Connect(ctx); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	defer mongodb.Disconnect(ctx)

	// kafka connsumer function
	kafkaConsumer := kafka.NewKafkaConsumer(mongodb, &loggs)

	// Kafka configuration
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	brokers := []string{"kafkamongo:9092"}
	groupID := "ledger-consumtion-group"
	topicName := []string{"transaction-ledger"}
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
			err = consumerGroup.Consume(ctx, topicName, kafkaConsumer)
			if err != nil {
				log.Printf("Consumer error: %v", err)
			}
			if ctx.Err() != nil {
				return
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
