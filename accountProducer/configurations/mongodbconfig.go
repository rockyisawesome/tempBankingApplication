package configurations

import "github.com/nicholasjackson/env" // Importing env package for environment variable parsing

// MongoDbConfig holds configuration details for connecting to a MongoDB instance.
// It encapsulates the MongoDB connection URI and the database name, which are typically
// sourced from environment variables for flexibility and security.
type MongoDbConfig struct {
	MongoURI string // MongoURI is the connection string for the MongoDB server (e.g., mongodb://user:pass@host:port/db?options)
	DBName   string // DBName is the name of the database to use within MongoDB
}

// NewMongoDbConfig creates a new MongoDbConfig instance by parsing environment variables.
// It uses the env package to retrieve the MongoDB URI and database name, providing default
// values if the environment variables are not set. Returns a pointer to the MongoDbConfig
// struct and an error if parsing fails.
func NewMongoDbConfig() (*MongoDbConfig, error) {
	// Define environment variable for MongoDB URI with a default value
	// The URI includes credentials, host, port, database, and authentication source
	var mongouri *string = env.String("MONGO_URI", false, "mongodb://admin:abcd@mongo:27017/ledger?authSource=admin", "Bind address for the Mongo server")

	// Define environment variable for database name with a default value
	var dbname *string = env.String("DB_NAME", false, "ledger", "Bind address for the Mongo server")

	// Parse environment variables; returns an error if parsing fails (e.g., invalid format)
	if err := env.Parse(); err != nil {
		return nil, err
	}

	// Construct and return the MongoDbConfig struct with parsed values
	return &MongoDbConfig{
		MongoURI: *mongouri, // Dereference the pointer to get the string value
		DBName:   *dbname,   // Dereference the pointer to get the string value
	}, nil
}
