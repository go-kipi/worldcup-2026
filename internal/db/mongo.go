package db

import (
	"context"
	"log"
	"time"

	"github.com/go-kipi/worldcup-2026/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoDatabase(cfg *config.Config) (*mongo.Database, error) {
	if cfg.MongoURI == "" {
		return nil, nil // Or handle as error if required
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI).
		SetBSONOptions(&options.BSONOptions{
			ObjectIDAsHexString: true,
		}))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Connected to MongoDB at %s", cfg.MongoURI)
	return client.Database(cfg.MongoDBName), nil
}

// AutoMigrate is a placeholder for MongoDB. In NoSQL, we don't usually migrate schema,
// but we can ensure indexes here.
func AutoMigrate(db *mongo.Database) error {
	// For now, just ensure we can reach the DB.
	// You can add index creation logic here for models like User (email unique index).
	log.Printf("MongoDB AutoMigrate: ensuring collections/indexes...")
	return nil
}
