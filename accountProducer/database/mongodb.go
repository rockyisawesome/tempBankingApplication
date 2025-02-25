package database

import (
	configs "accountProducer/configurations" // Importing configurations package for MongoDB settings
	"accountProducer/models"                 // Importing models package for TransactionLedger struct
	"context"                                // Importing context for request-scoped operations and cancellation
	"fmt"                                    // Importing fmt for error formatting
	"time"                                   // Importing time for setting query timeouts

	"github.com/hashicorp/go-hclog"                // Importing hclog for structured logging
	"go.mongodb.org/mongo-driver/v2/bson"          // Importing bson for MongoDB query construction
	"go.mongodb.org/mongo-driver/v2/mongo"         // Importing mongo for MongoDB client operations
	"go.mongodb.org/mongo-driver/v2/mongo/options" // Importing options for MongoDB client configuration
)

// MongoDB represents a MongoDB database connection.
// It implements the Database interface and encapsulates the MongoDB client, database instance,
// configuration, and logger for interacting with MongoDB.
type MongoDB struct {
	Client   *mongo.Client          // Client is the MongoDB client instance for database operations
	Database *mongo.Database        // Database is the specific MongoDB database instance
	Config   *configs.MongoDbConfig // Config holds the MongoDB connection settings (e.g., URI, DB name)
	loggs    *hclog.Logger          // loggs is the logger instance for logging database operations
	// Ctx      *context.Context    // Commented out: Context is passed per method, not stored
}

// NewMongoDB creates a new MongoDB instance with the provided configuration and logger.
// It initializes the struct but does not establish a connection yet; Connect must be called separately.
// Returns a pointer to the MongoDB struct.
func NewMongoDB(cfg *configs.MongoDbConfig, lobbs *hclog.Logger) *MongoDB {
	return &MongoDB{
		Config: cfg,   // Set the MongoDB configuration
		loggs:  lobbs, // Set the logger instance
	}
}

// Connect establishes a connection to the MongoDB database using the provided URI from Config.
// It pings the database to verify the connection and sets the Client and Database fields on success.
// Takes a context for cancellation and timeout support. Returns an error if connection or ping fails.
func (mango *MongoDB) Connect(ctx context.Context) error {
	// Create a new MongoDB client with the configured URI
	client, err := mongo.Connect(options.Client().ApplyURI(mango.Config.MongoURI))
	if err != nil {
		(*mango.loggs).Error("Error connecting to Mongo DB", "Error", err)
		return err
	}

	// Ping the database to ensure it's reachable
	err = client.Ping(ctx, nil)
	if err != nil {
		(*mango.loggs).Error("Pinging Database but no response", "Error", err)
		return err
	}
	(*mango.loggs).Info("Database is up and active", "Error", err)

	// Store the client and database instances in the struct
	mango.Client = client
	mango.Database = client.Database(mango.Config.DBName)
	(*mango.loggs).Info("Connected to Database", "DB", mango.Config.DBName)

	return nil
}

// Disconnect closes the connection to the MongoDB database.
// Implements the Database interface's Disconnect method. Uses the provided context for cancellation
// and timeout handling. Returns an error if disconnection fails (e.g., due to network issues).
func (mango *MongoDB) Disconnect(ctx context.Context) error {
	// Disconnect the MongoDB client
	return mango.Client.Disconnect(ctx)
}

// GetTransactionsByAccountNumber retrieves all transactions for a given account number from MongoDB.
// Implements the Database interface's method. Queries the "transactions" collection using the
// accountNumber as a filter on the "from_account_id" field. Returns a pointer to a slice of
// TransactionLedger structs or an error if the query or decoding fails.
func (mango *MongoDB) GetTransactionsByAccountNumber(ctx context.Context, accountNumber string) (*[]models.TransactionLedger, error) {
	// Check if the database connection is initialized
	if mango.Database == nil {
		return nil, fmt.Errorf("database not initialized, call Connect first")
	}

	// Access the "transactions" collection
	collection := mango.Database.Collection("transactions")

	// Set a 10-second timeout for the query to prevent hanging
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() // Ensure the timeout is cleaned up

	// Define the filter to match transactions by account number
	filter := bson.M{"from_account_id": accountNumber}

	// Execute the find query to retrieve matching documents
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		(*mango.loggs).Error("Failed to query transactions", "accountNumber", accountNumber, "Error", err)
		return nil, fmt.Errorf("failed to query transactions for %s: %w", accountNumber, err)
	}
	defer cursor.Close(ctx) // Ensure the cursor is closed after use

	// Decode all matching documents into a slice of TransactionLedger
	var transactions []models.TransactionLedger
	if err = cursor.All(ctx, &transactions); err != nil {
		(*mango.loggs).Error("Failed to decode transactions", "accountNumber", accountNumber, "Error", err)
		return nil, fmt.Errorf("failed to decode transactions for %s: %w", accountNumber, err)
	}

	// Log successful retrieval with the number of transactions found
	(*mango.loggs).Info("Successfully retrieved transactions", "accountNumber", accountNumber, "Count", len(transactions))

	// Return the slice of transactions
	return &transactions, nil
}
