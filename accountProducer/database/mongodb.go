package database

import (
	configs "accountProducer/configurations"
	"accountProducer/models"
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	Config   *configs.MongoDbConfig
	loggs    *hclog.Logger
	// Ctx      *context.Context
}

// returning an instance og MongoDB
func NewMongoDB(cfg *configs.MongoDbConfig, lobbs *hclog.Logger) *MongoDB {
	return &MongoDB{
		Config: cfg,
		loggs:  lobbs,
	}
}

func (mango *MongoDB) Connect(ctx context.Context) error {

	client, err := mongo.Connect(options.Client().ApplyURI(mango.Config.MongoURI))
	if err != nil {
		(*mango.loggs).Error("Error connecting to Mongo DB", "Error", err)
		return err
	}

	// ping the database to verify connections
	err = client.Ping(ctx, nil)
	if err != nil {
		(*mango.loggs).Error("Pinging Database but no response", "Error", err)
		return err
	}
	(*mango.loggs).Info("Database is up and active", "Error", err)

	mango.Client = client
	mango.Database = client.Database(mango.Config.DBName)
	(*mango.loggs).Info("Connected to Database", "DB", mango.Config.DBName)

	return nil
}

// Disconnect implements the Database interface
func (mango *MongoDB) Disconnect(ctx context.Context) error {

	return mango.Client.Disconnect(ctx)

}

// GetTransactionsByAccountNumber retrieves transactions for a given account number
func (mango *MongoDB) GetTransactionsByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error) {
	// Ensure the database connection is established
	if mango.Database == nil {
		return nil, fmt.Errorf("database not initialized, call Connect first")
	}

	// Define the collection
	collection := mango.Database.Collection("transactions")

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Define the filter to match accountNumber
	filter := bson.M{"accountNumber": accountNumber}

	// Execute the query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		(*mango.loggs).Error("Failed to query transactions", "accountNumber", accountNumber, "Error", err)
		return nil, fmt.Errorf("failed to query transactions for %s: %w", accountNumber, err)
	}
	defer cursor.Close(ctx)

	// Decode results into a slice
	var transactions []models.TransactionLedger
	if err = cursor.All(ctx, &transactions); err != nil {
		(*mango.loggs).Error("Failed to decode transactions", "accountNumber", accountNumber, "Error", err)
		return nil, fmt.Errorf("failed to decode transactions for %s: %w", accountNumber, err)
	}

	// Log success
	(*mango.loggs).Info("Successfully retrieved transactions", "accountNumber", accountNumber, "Count", len(transactions))

	// Return the transactions
	return &transactions, nil
}
