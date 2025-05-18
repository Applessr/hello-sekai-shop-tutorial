package migration

import (
	"context"
	"log"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func inventoryDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConnect(pctx, cfg).Database("inventory_db")
}

func InventoryMigrate(pctx context.Context, cfg *config.Config) {
	db := inventoryDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("players_inventory")

	index, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"_id", 1}, {"item_id", 1}}},
	})
	for _, index := range index {
		log.Printf("index: %s", index)
	}

	col = db.Collection("players_inventory_queue")

	results, err := col.InsertOne(pctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate inventory completed: ", results)
}
