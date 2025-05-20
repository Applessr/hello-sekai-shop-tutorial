package migration

import (
	"context"
	"log"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func authDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("auth_db")
}

func AuthMigrate(pctx context.Context, cfg *config.Config) {
	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("auth")

	index, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"player_id", 1}}},
		{Keys: bson.D{{"refresh_token", 1}}},
	})
	for _, index := range index {
		log.Printf("index: %s", index)
	}

	col = db.Collection("roles")

	index, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"code", 1}}},
	})
	for _, index := range index {
		log.Printf("index: %s", index)
	}

	//roles
	documents := func() []any {
		roles := []*auth.Role{
			{
				Title: "player",
				Code:  0,
			},
			{
				Title: "admin",
				Code:  1,
			},
		}

		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, documents, nil)
	if err != nil {
		log.Fatalf("Error: InsertMany failed: %v", err)
	}
	log.Printf("Inserted %d documents into roles collection", len(results.InsertedIDs))
}
