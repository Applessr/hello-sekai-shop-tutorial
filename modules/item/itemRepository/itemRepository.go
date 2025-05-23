package itemRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ItemRepositoryService interface {
		InsertOneItem(pctx context.Context, req *item.Item) (primitive.ObjectID, error)
		IsUniqueItem(pctx context.Context, title string) bool
		FindOneItem(pctx context.Context, itemId string) (*item.Item, error)
		FindManyItems(pctx context.Context, filter primitive.D, option []*options.FindOptions) ([]*item.ItemShowCase, error)
		CountItems(pctx context.Context, filter primitive.D) (int64, error)
		UpdateOneItem(pctx context.Context, itemId string, req primitive.M) error
		EnableOrDisableItem(pctx context.Context, itemId string, isActive bool) error
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{db}
}

func (r *itemRepository) itemDbConnect(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}

func (r *itemRepository) InsertOneItem(pctx context.Context, req *item.Item) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	itemId, err := col.InsertOne(ctx, req, nil)
	if err != nil {
		log.Printf("Error: InsertOneItem: %s", err.Error())
		return primitive.ObjectID{}, errors.New("error: Insert one item failed")
	}
	return itemId.InsertedID.(primitive.ObjectID), nil
}

func (r *itemRepository) IsUniqueItem(pctx context.Context, title string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	result := new(item.Item)
	if err := col.FindOne(
		ctx,
		bson.M{"title": title}).Decode(result); err != nil {
		log.Printf("Error: IsUniqueItem: %s", err.Error())
		return true
	}

	return false
}

func (r *itemRepository) FindOneItem(pctx context.Context, itemId string) (*item.Item, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	result := new(item.Item)
	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}).Decode(result); err != nil {
		log.Printf("Error: FindOneItem: %s", err.Error())
		return nil, errors.New("error: item not found")
	}

	return result, nil
}

func (r *itemRepository) FindManyItems(pctx context.Context, filter primitive.D, option []*options.FindOptions) ([]*item.ItemShowCase, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	cursors, err := col.Find(ctx, filter, option...)
	if err != nil {
		log.Printf("Error: FindManyItems: %s", err.Error())
		return nil, errors.New("error: find many items failed")
	}

	results := make([]*item.ItemShowCase, 0)
	for cursors.Next(ctx) {
		result := new(item.Item)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: FindManyItems: %s", err.Error())
			return make([]*item.ItemShowCase, 0), errors.New("error: find many items failed")
		}
		results = append(results, &item.ItemShowCase{
			ItemId:   "item:" + result.Id.Hex(),
			Title:    result.Title,
			Price:    result.Price,
			Damage:   result.Damage,
			ImageUrl: result.ImageUrl,
		})
	}

	return results, nil
}

func (r *itemRepository) CountItems(pctx context.Context, filter primitive.D) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	count, err := col.CountDocuments(ctx, filter)
	if err != nil && count == 0 {
		log.Printf("Error: CountItems: %s", err.Error())
		return -1, errors.New("error: count items failed")
	}

	return count, nil
}

func (r *itemRepository) UpdateOneItem(pctx context.Context, itemId string, req primitive.M) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}, bson.M{"$set": req})
	if err != nil {
		log.Printf("Error: UpdateOneItem failed: %s", err.Error())
		return errors.New("error: update one item failed")
	}
	log.Printf("UpdateOneItem result: %v", result.ModifiedCount)

	return nil
}

func (r *itemRepository) EnableOrDisableItem(pctx context.Context, itemId string, isActive bool) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConnect(ctx)
	col := db.Collection("item")

	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvertToObjectId(itemId)}, bson.M{"$set": bson.M{"usage_status": isActive}})
	if err != nil {
		log.Printf("Error: EnableOrDisableItem failed: %s", err.Error())
		return errors.New("error: enable or disable item failed")
	}
	log.Printf("EnableOrDisableItem result: %v", result.ModifiedCount)

	return nil
}
