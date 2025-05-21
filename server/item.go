package server

import (
	"log"

	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemHandler"
	itemPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemUsecase"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/grpccon"
)

func (s *server) itemService() {
	repo := itemRepository.NewItemRepository(s.db)
	usecase := itemUsecase.NewItemUsecase(repo)
	httpHandler := itemHandler.NewItemHttpHandler(s.cfg, usecase)
	grpcHandler := itemHandler.NewItemGrpcHandler(usecase)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemPb.RegisterItemGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Item gRPC server listening on %s", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(lis)
	}()

	_ = grpcHandler

	item := s.app.Group("/item_v1")

	item.GET("", s.healthCheckService)
	item.GET("/item", httpHandler.FindManyItem)
	item.GET("/item/:item_id", httpHandler.FindOneItem)

	item.POST("/item", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(httpHandler.CreatedItem, []int{1, 0})))

	item.PATCH("/item/:item_id", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(httpHandler.EditItem, []int{1, 0})))
	item.PATCH("/item/:item_id/is-activated", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(httpHandler.EnableOrDisableItem, []int{1, 0})))
}
