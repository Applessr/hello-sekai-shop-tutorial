package database

import (
	"context"
	"log"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbConnect(pctx context.Context, cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.Url))
	if err != nil {
		log.Fatal("Error:Connect to database failed:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Error: Pinging to database failed:", err)
	}

	return client

}
