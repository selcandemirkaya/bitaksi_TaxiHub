package repository

//mongo bağlantısı

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//secret management tool (hashicorp vault)--> büyük projeler için

type Mongo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func ConnectDatabase() *Mongo {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "taxiHub"
	}
	colName := os.Getenv("MONGO_COLLECTION")
	if colName == "" {
		colName = "drivers"
	}

	//timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//client
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Mongo connect error: %v", err)
	}

	// Bağlantıyı test et (ping)
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Mongo ping error: %v", err)
	}

	collection := client.Database(dbName).Collection(colName)

	log.Printf("Mongo connected. DB=%s Collection=%s", dbName, colName)

	return &Mongo{
		Client:     client,
		Collection: collection,
	}

}

func (m *Mongo) CloseMongo() {
	if m.Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		log.Printf("Mongo disconnect error: %v", err)
		return
	}

	log.Println("Mongo disconnected.")
}
