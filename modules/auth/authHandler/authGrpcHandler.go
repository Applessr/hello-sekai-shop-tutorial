package authHandler

import "github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authUsecase"

type (
	authGrpcHandler struct {
		authUsecase authUsecase.AuthUsecaseService
	}
)

func NewAuthGrpcHandler(authUsecase authUsecase.AuthUsecaseService) *authGrpcHandler {
	return &authGrpcHandler{authUsecase}
}
