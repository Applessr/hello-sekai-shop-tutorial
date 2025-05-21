package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	authPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/grpccon"
)

type (
	MiddlewareRepositoryService interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
	}

	middlewareRepository struct{}
)

func NewMiddlewareRepository() MiddlewareRepositoryService {
	return &middlewareRepository{}
}

func (r *middlewareRepository) AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRpc client connection failed: %s", err.Error())
		return errors.New("error: gRpc client connection failed")
	}

	result, err := conn.Auth().AccessTokenSearch(ctx, &authPb.AccessTokenSearchReq{AccessToken: accessToken})
	if err != nil {
		log.Printf("Error: AccessTokenSearch: %s", err.Error())
		return errors.New(err.Error())
	}

	if result == nil {
		return errors.New("error: access token not found")
	}

	if !result.IsValid {
		return errors.New("error: access token is invalid")
	}

	return nil
}
