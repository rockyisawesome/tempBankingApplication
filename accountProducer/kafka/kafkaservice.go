package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaController struct {
}

// registering the routes
func (k *KafkaController) PushOrderToQueue(topic string, message []byte) error {

	brokers := []string{"kafka:9092"}
	producer, err := k.connectProducer(brokers, topic)
	if err != nil {
		return err
	}

	defer producer.Close()

	// create new kafka message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// send message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	fmt.Println("Account post request is stored", topic, partition, offset)
	return nil

}

func (k *KafkaController) connectProducer(brokers []string, topic string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create cluster admin: %v", err)
	}
	defer admin.Close()

	// Specify topic details
	partitions := int32(3)        // Number of partitions
	replicationFactor := int16(1) // Replication factor (1 for single broker)

	// Create topic if it doesn't exist
	err = k.createTopicIfNotExists(admin, topic, partitions, replicationFactor)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// time.Sleep(2 * time.Second)

	return sarama.NewSyncProducer(brokers, config)
}

// createTopicIfNotExists creates a Kafka topic if it doesn't already exist.
func (k *KafkaController) createTopicIfNotExists(admin sarama.ClusterAdmin, topicName string, partitions int32, replicationFactor int16) error {
	// List existing topics
	topics, err := admin.ListTopics()
	fmt.Println(topics)
	if err != nil {
		fmt.Errorf("failed to list topics: %v", err)
	}

	// Check if the topic already exists
	if _, exists := topics[topicName]; exists {
		fmt.Printf("Topic %s already exists, skipping creation\n", topicName)
		return nil
	}

	// Define topic details
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries: map[string]*string{
			"retention.ms": k.stringPtr("604800000"), // 7 days retention (optional)
		},
	}

	// Create the topic
	err = admin.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %v", topicName, err)
	}

	fmt.Printf("Topic %s created successfully\n", topicName)
	time.Sleep(2 * time.Second)
	return nil
}

// Helper function to convert string to *string
func (k *KafkaController) stringPtr(s string) *string {
	return &s
}
