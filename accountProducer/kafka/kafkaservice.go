package kafka

import (
	"fmt"  // Importing fmt for formatted output and error construction
	"log"  // Importing log for fatal error logging
	"time" // Importing time for delays in topic creation

	"github.com/IBM/sarama" // Importing sarama for Kafka client functionality
)

// KafkaController manages interactions with a Kafka cluster.
// It provides methods to push messages to Kafka topics and handle topic creation.
type KafkaController struct {
	// No fields are currently needed, but this struct could be extended for configuration or state
}

// PushOrderToQueue sends a message to a specified Kafka topic.
// It connects to a Kafka producer, constructs a message, and sends it synchronously.
// Takes a topic name and a byte slice message as input. Returns an error if the operation fails.
func (k *KafkaController) PushOrderToQueue(topic string, message []byte) error {
	// Define the Kafka broker address (hardcoded for simplicity)
	brokers := []string{"kafka:9092"}

	// Establish a connection to the Kafka producer
	producer, err := k.connectProducer(brokers, topic)
	if err != nil {
		return err // Return the error if producer initialization fails
	}
	defer producer.Close() // Ensure the producer is closed after use

	// Create a new Kafka message with the provided topic and message content
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message), // Encode the byte slice as a string
	}

	// Send the message synchronously and retrieve partition and offset
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err // Return the error if sending fails
	}

	// Log success with topic, partition, and offset details
	fmt.Println("Account post request is stored", topic, partition, offset)
	return nil
}

// connectProducer initializes a synchronous Kafka producer and ensures the topic exists.
// Takes a list of broker addresses and the topic name as input. Returns a SyncProducer
// instance or an error if the connection or topic creation fails.
func (k *KafkaController) connectProducer(brokers []string, topic string) (sarama.SyncProducer, error) {
	// Create a new Kafka configuration
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true          // Ensure successful sends are reported
	config.Producer.RequiredAcks = sarama.WaitForAll // Require acknowledgment from all replicas
	config.Producer.Retry.Max = 5                    // Set maximum retries for failed sends

	// Create a cluster admin client to manage topics
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create cluster admin: %v", err) // Fatal error if admin creation fails
	}
	defer admin.Close() // Ensure the admin client is closed after use

	// Specify topic details
	partitions := int32(3)        // Number of partitions for the topic
	replicationFactor := int16(1) // Replication factor (1 for a single broker setup)

	// Create the topic if it doesn’t already exist
	err = k.createTopicIfNotExists(admin, topic, partitions, replicationFactor)
	if err != nil {
		log.Fatalf("Error: %v", err) // Fatal error if topic creation fails
	}

	// Initialize and return a synchronous producer
	return sarama.NewSyncProducer(brokers, config)
}

// createTopicIfNotExists creates a Kafka topic if it doesn’t already exist.
// Takes a cluster admin client, topic name, number of partitions, and replication factor as input.
// Returns an error if listing topics or creating the topic fails.
func (k *KafkaController) createTopicIfNotExists(admin sarama.ClusterAdmin, topicName string, partitions int32, replicationFactor int16) error {
	// Retrieve the list of existing topics
	topics, err := admin.ListTopics()
	fmt.Println(topics) // Log the current topics for debugging
	if err != nil {
		return fmt.Errorf("failed to list topics: %v", err)
	}

	// Check if the topic already exists to avoid redundant creation
	if _, exists := topics[topicName]; exists {
		fmt.Printf("Topic %s already exists, skipping creation\n", topicName)
		return nil
	}

	// Define the topic configuration
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
		ConfigEntries: map[string]*string{
			"retention.ms": k.stringPtr("604800000"), // Set retention to 7 days (604800000 ms)
		},
	}

	// Create the topic with the specified details
	err = admin.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %v", topicName, err)
	}

	// Log success and introduce a brief delay to ensure topic availability
	fmt.Printf("Topic %s created successfully\n", topicName)
	time.Sleep(2 * time.Second) // Wait for topic propagation
	return nil
}

// stringPtr is a helper function that converts a string to a pointer to a string.
// Used for setting configuration entries in topic creation. Takes a string as input
// and returns a pointer to it.
func (k *KafkaController) stringPtr(s string) *string {
	return &s
}
