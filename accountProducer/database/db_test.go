package database

import (
	"context"
	"testing"

	configs "accountProducer/configurations"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MockMongoClient is a mock for mongo.Client
type MockMongoClient struct {
	*mongo.Client
	PingFunc       func(ctx context.Context, rp *readpref.ReadPref) error
	DisconnectFunc func(ctx context.Context) error
}

func (m *MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.PingFunc(ctx, rp)
}

func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	return m.DisconnectFunc(ctx)
}

// TestNewMongoDB tests the constructor
func TestNewMongoDB(t *testing.T) {
	cfg := &configs.MongoDbConfig{
		MongoURI: "mongodb://admin:abcd@mongo:27017/ledger?authSource=admin",
		DBName:   "ledger",
	}
	logger := hclog.NewNullLogger()
	db := NewMongoDB(cfg, &logger)

	assert.NotNil(t, db)
	assert.Equal(t, cfg, db.Config)
	assert.Equal(t, &logger, db.loggs)
	assert.Nil(t, db.Client)
	assert.Nil(t, db.Database)
}
