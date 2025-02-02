package utils

import (
	"context"
	"hollow/configs"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(configs.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(configs.DBName)
	log.Println("Connected to MongoDB!")
}
