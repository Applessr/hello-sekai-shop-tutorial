package inventoryRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	InventoryRepositoryService interface{}

	InventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepositoryService {
	return &InventoryRepository{db}
}

func (r *InventoryRepository) inventoryDbConnect(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}
