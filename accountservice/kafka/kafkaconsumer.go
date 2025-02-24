package kafka

import (
	"accountservice/database"
	"accountservice/models"
	"accountservice/repositories"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	repo repositories.Repository
}

func NewKafkaConsumer(db *database.PostgresPoolDB) *KafkaConsumer {
	return &KafkaConsumer{
		repo: repositories.NewUserRepository(db),
	}
}

func (h KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var account models.Account
		err := json.Unmarshal(msg.Value, &account)
		if err != nil {
			log.Printf("Failed to unmarshal message (offset %d): %v", msg.Offset, err)
			continue
		}

		ctx := context.Background()
		// Process the account (e.g., save to DB)
		err = h.repo.Create(ctx, &account)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		log.Printf("Processed account: %+v (partition %d, offset %d)", account, msg.Partition, msg.Offset)
		fmt.Println("Account is saved in postgres")
		// Mark message as processed (commit offset)
		session.MarkMessage(msg, "")
	}

	return nil
}

func (KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
