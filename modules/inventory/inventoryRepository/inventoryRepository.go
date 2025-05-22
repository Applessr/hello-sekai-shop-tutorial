package inventoryRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/modules/inventory"
	itemPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/models"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/grpccon"
	jwtAuth "github.com/Applessr/hello-sekai-shop-tutorial/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryRepositoryService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		FindItemInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemInIdsReq) (*itemPb.FindItemInIdsRes, error)
		FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error)
		CountPlayerItems(pctx context.Context, playerId string) (int64, error)
	}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepositoryService {
	return &inventoryRepository{db}
}

func (r *inventoryRepository) inventoryDbConnect(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}

func (r *inventoryRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConnect(ctx)
	col := db.Collection("players_inventory_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("Error: GetOffset failed: %s", err.Error())
		return -1, errors.New("error: GetOffset failed")
	}

	return result.Offset, nil
}

func (r *inventoryRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConnect(ctx)
	col := db.Collection("players_inventory_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf("Error: UpserOffset failed: %s", err.Error())
		return errors.New("error: UpserOffset failed")
	}
	log.Printf("Info: UpserOffset result: %v", result)

	return nil
}

func (r *inventoryRepository) FindItemInIds(pctx context.Context, grpcUrl string, req *itemPb.FindItemInIdsReq) (*itemPb.FindItemInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtAuth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRpc client connection failed: %s", err.Error())
		return nil, errors.New("error: gRpc client connection failed")
	}

	result, err := conn.Item().FindItemInIds(ctx, req)
	if err != nil {
		log.Printf("Error: FindItemInIds: %s", err.Error())
		return nil, errors.New(err.Error())
	}

	if result == nil && len(result.Items) == 0 {
		return nil, errors.New("error: item not found")
	}

	return result, nil
}

func (r *inventoryRepository) FindPlayerItems(pctx context.Context, filter primitive.D, opts []*options.FindOptions) ([]*inventory.Inventory, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConnect(ctx)
	col := db.Collection("players_inventory")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindPlayerItems failed: %s", err.Error())
		return nil, errors.New("error: player items not found")
	}

	results := make([]*inventory.Inventory, 0)
	for cursors.Next(ctx) {
		result := new(inventory.Inventory)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: FindPlayerItems failed: %s", err.Error())
			return nil, errors.New("error: player items not found")
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *inventoryRepository) CountPlayerItems(pctx context.Context, playerId string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConnect(ctx)
	col := db.Collection("players_inventory")

	count, err := col.CountDocuments(ctx, bson.M{"player_id": playerId})
	if err != nil {
		log.Printf("Error: CountPlayerItems failed: %s", err.Error())
		return -1, errors.New("error: count player items failed")
	}

	return count, nil
}
