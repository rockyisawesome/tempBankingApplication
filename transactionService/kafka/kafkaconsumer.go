package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"transactionService/database"
	"transactionService/models"
	"transactionService/repositories"

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
		var trans models.Transaction
		err := json.Unmarshal(msg.Value, &trans)
		if err != nil {
			log.Printf("Failed to unmarshal message (offset %d): %v", msg.Offset, err)
			continue
		}

		ctx := context.Background()
		fmt.Println(trans)
		// Process the transaction
		kafkapush := KafkaController{}
		err = h.repo.TransactionRouter(ctx, &trans)
		if err != nil {
			fmt.Println(err)
			//push to dead order queue
			err = kafkapush.PushToQueue("dead-ledger", &trans)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		}
		// push transacation into mongo ledger service
		err = kafkapush.PushToQueue("transaction-ledger", &trans)
		if err != nil {
			fmt.Println(err)
			//push to dead order queue
			err = kafkapush.PushToQueue("dead-ledger", &trans)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return nil
		}

		log.Printf("Processed trnsaction: %+v (partition %d, offset %d)", trans.FromAccountID, msg.Partition, msg.Offset)
		fmt.Println("Account is saved in postgres")
		// Mark message as processed (commit offset)
		session.MarkMessage(msg, "")
	}

	return nil
}

func (KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
