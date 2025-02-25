package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"ledgerservice/database"
	"ledgerservice/models"
	"ledgerservice/repositories"
	"log"

	"github.com/IBM/sarama"
	"github.com/hashicorp/go-hclog"
)

type KafkaConsumer struct {
	repo  repositories.Repository
	loggs *hclog.Logger
}

func NewKafkaConsumer(db database.Database, lobbs *hclog.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		repo:  repositories.NewTransactionRepository(db, lobbs),
		loggs: lobbs,
	}
}

func (h KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var trans models.TransactionLedger
		err := json.Unmarshal(msg.Value, &trans)
		if err != nil {
			log.Printf("Failed to unmarshal message (offset %d): %v", msg.Offset, err)
			continue
		}

		ctx := context.Background()
		fmt.Println(trans)
		// Add the transaction into transaction ledger
		err = h.repo.InsertTransaction(ctx, trans)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		log.Printf("Processed Transaction: %+v (partition %d, offset %d)", trans, msg.Partition, msg.Offset)
		fmt.Println("Transaction is logged")
		// Mark message as processed (commit offset)
		session.MarkMessage(msg, "")
	}

	return nil
}

func (KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
