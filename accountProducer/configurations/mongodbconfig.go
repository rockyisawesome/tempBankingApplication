package configurations

import "github.com/nicholasjackson/env"

type MongoDbConfig struct {
	MongoURI string
	DBName   string
}

// mongodb://admin:secret@mongo:27017/ledger?authSource=admin

func NewMongoDbConfig() (*MongoDbConfig, error) {
	var mongouri *string = env.String("MONGO_URI", false, "mongodb://admin:abcd@localhost:27017/ledger?authSource=admin", "Bind address for the Mongo server")
	var dbname *string = env.String("DB_NAME", false, "ledger", "Bind address for the Mongo server")
	if err := env.Parse(); err != nil {
		return nil, err
	}

	return &MongoDbConfig{
		MongoURI: *mongouri,
		DBName:   *dbname,
	}, nil
}
