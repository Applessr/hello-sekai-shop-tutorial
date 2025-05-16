package inventoryHandler

import "github.com/Applessr/hello-sekai-shop-tutorial/modules/inventory/inventoryUsecase"

type (
	inventoryGrpcHandler struct {
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{inventoryUsecase}
}
