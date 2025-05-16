package middlewareUsecase

import "github.com/Applessr/hello-sekai-shop-tutorial/modules/middleware/middlewareRepository"

type (
	MiddlewareUsecaseService interface{}

	middlewareUsecase struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{middlewareRepository}
}
