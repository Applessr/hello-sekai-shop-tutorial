package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	authPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/grpccon"
	jwtAuth "github.com/Applessr/hello-sekai-shop-tutorial/pkg/jwtauth"
)

type (
	MiddlewareRepositoryService interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
		RolesCount(pctx context.Context, grpcUrl string) (int64, error)
	}

	middlewareRepository struct{}
)

func NewMiddlewareRepository() MiddlewareRepositoryService {
	return &middlewareRepository{}
}

func (r *middlewareRepository) AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	jwtAuth.SetApiKeyInContext(&ctx)
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
		log.Printf("Error: access token not found")
		return errors.New("error: access token not found")
	}

	if !result.IsValid {
		log.Printf("Error: access token is invalid")
		return errors.New("error: access token is invalid")
	}

	return nil
}

func (r *middlewareRepository) RolesCount(pctx context.Context, grpcUrl string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	jwtAuth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRpc client connection failed: %s", err.Error())
		return -1, errors.New("error: gRpc client connection failed")
	}

	result, err := conn.Auth().RolesCount(ctx, &authPb.RolesCountReq{})
	if err != nil {
		log.Printf("Error: RolesCount: %s", err.Error())
		return -1, errors.New(err.Error())
	}

	if result == nil {
		log.Printf("Error: roles count failed")
		return -1, errors.New("error: roles count failed")
	}

	return result.Count, nil
}
