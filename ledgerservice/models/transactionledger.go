package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionLedger struct {
	ID              bson.ObjectID `bson:"_id"`
	FromAccountID   string        `bson:"from_account_id" json:"from_account_id"`
	ToAccountID     string        `bson:"to_account_id" json:"to_account_id"`
	Amount          float64       `bson:"amount" json:"amount"`
	TransactionType string        `bson:"transaction_type" json:"transaction_type"`
	Description     string        `bson:"description" json:"description"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	Status          string        `bson:"status" json:"status"`
}
