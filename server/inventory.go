package server

import (
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/inventory/inventoryHandler"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/inventory/inventoryRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/inventory/inventoryUsecase"
)

func (s *server) inventoryService() {
	repo := inventoryRepository.NewInventoryRepository(s.db)
	usecase := inventoryUsecase.NewInventoryUsecase(repo)
	httpHandler := inventoryHandler.NewInventoryHttpHandler(s.cfg, usecase)
	queueHandler := inventoryHandler.NewInventoryQueueHandler(s.cfg, usecase)

	go queueHandler.AddPlayerItem()
	go queueHandler.RollbackAddPlayerItem()
	go queueHandler.RemovePlayerItem()
	go queueHandler.RollbackRemovePlayerItem()

	inventory := s.app.Group("/inventory_v1")

	inventory.GET("", s.healthCheckService)
	inventory.GET("/inventory/:player_id", httpHandler.FindPlayerItems, s.middleware.JwtAuthorization, s.middleware.PlayerIdParamValidation)
}
