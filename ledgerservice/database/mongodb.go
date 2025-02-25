package database

import (
	"context"
	configs "ledgerservice/configurations"
	"ledgerservice/models"
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

// InsertTransaction inserts a TransactionLedger document into the transactions collection
func (mango *MongoDB) InsertTransaction(ctx context.Context, ledger models.TransactionLedger) (string, error) {
	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Access the 'transactions' collection
	collection := mango.Database.Collection("transactions")

	// Insert the document
	result, err := collection.InsertOne(ctx, ledger)
	if err != nil {
		(*mango.loggs).Error("Failed to insert transaction into MongoDB", "Error", err)
		return "", err
	}

	// Get the inserted ID and convert it to a string
	insertedID := result.InsertedID.(bson.ObjectID).Hex()
	(*mango.loggs).Info("Successfully inserted transaction", "ID", insertedID)

	return insertedID, nil
}
