package server

import (
	"log"

	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authHandler"
	authPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authUsecase"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/grpccon"
)

func (s *server) authService() {
	repo := authRepository.NewAuthRepository(s.db)
	usecase := authUsecase.NewAuthUsecase(repo)
	httpHandler := authHandler.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandler.NewAuthGrpcHandler(usecase)

	//gRPC
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()

	auth := s.app.Group("/auth_v1")

	// auth.GET("", s.middleware.JwtAuthorization(s.middleware.RbacAuthorization(s.healthCheckService, []int{1, 0})))
	// auth.GET("/test/:player_id", s.healthCheckService, s.middleware.JwtAuthorization, s.middleware.PlayerIdParamValidation)
	auth.GET("", s.healthCheckService)

	auth.POST("/auth/login", httpHandler.Login)
	auth.POST("/auth/refresh-token", httpHandler.RefreshToken)
	auth.POST("/auth/logout", httpHandler.Logout)
}
