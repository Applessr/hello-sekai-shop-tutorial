package migration

import (
	"context"
	"log"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/player"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/database"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func playerDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("player_db")
}

func PlayerMigrate(pctx context.Context, cfg *config.Config) {
	db := playerDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("player_transactions")

	index, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"player_id", 1}}},
	})
	for _, index := range index {
		log.Printf("index: %s", index)
	}

	col = db.Collection("players")

	index, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}}},
		{Keys: bson.D{{"email", 1}}},
	})
	for _, index := range index {
		log.Printf("index: %s", index)
	}

	documents := func() []any {
		roles := []*player.Player{
			{
				Email:    "player001@mail.com",
				Password: "123456",
				Username: "Player001",
				PlayerRole: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:    "player002@mail.com",
				Password: "123456",
				Username: "Player002",
				PlayerRole: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:    "player003@mail.com",
				Password: "123456",
				Username: "Player003",
				PlayerRole: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:    "admin@mail.com",
				Password: "123456",
				Username: "Admin1",
				PlayerRole: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
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

	playerTransaction := make([]any, 0)
	for _, p := range results.InsertedIDs {
		playerTransaction = append(playerTransaction, &player.PlayerTransaction{
			PlayerId:  "player:" + p.(primitive.ObjectID).Hex(),
			Amount:    1000,
			CreatedAt: utils.LocalTime(),
		})

	}

	col = db.Collection("player_transactions")
	results, err = col.InsertMany(pctx, playerTransaction, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Inserted %d documents into player_transactions collection", len(results.InsertedIDs))

	col = db.Collection("player_transactions_queue")
	result, err := col.InsertOne(pctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Inserted %d documents into player_transactions_queue collection", result)
}
