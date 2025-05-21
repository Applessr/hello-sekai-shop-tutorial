package authUsecase

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth"
	authPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth/authRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/player"
	playerPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/player/playerPb"
	jwtAuth "github.com/Applessr/hello-sekai-shop-tutorial/pkg/jwtauth"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/utils"
)

type (
	AuthUsecaseService interface {
		Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileInterceptor, error)
		RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileInterceptor, error)
		Logout(pctx context.Context, credentialId string) (int64, error)
		AccessTokenSearch(pctx context.Context, accessToken string) (*authPb.AccessTokenSearchRes, error)
		RoleCount(pctx context.Context) (*authPb.RolesCountRes, error)
	}

	authUsecase struct {
		authRepository authRepository.AuthRepositoryService
	}
)

func NewAuthUsecase(authRepository authRepository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{authRepository}
}

func (u *authUsecase) Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileInterceptor, error) {
	profile, err := u.authRepository.CredentialSearch(pctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	profile.Id = "player:" + profile.Id

	accessToken := jwtAuth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtAuth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtAuth.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwtAuth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	credentialId, err := u.authRepository.InsertOneCredential(pctx, &auth.Credential{
		PlayerId:     profile.Id,
		RoleCode:     int(profile.RoleCode),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})

	credential, err := u.authRepository.FindOnePlayerCredential(pctx, credentialId.Hex())
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileInterceptor{
		PlayerProfile: &player.PlayerProfile{
			Id:        profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			PlayerId:     credential.PlayerId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUsecase) RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileInterceptor, error) {
	claims, err := jwtAuth.ParseToken(cfg.Jwt.RefreshSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: Parse token failed: %s", err)
		return nil, errors.New(err.Error())
	}

	profile, err := u.authRepository.FindOnePlayerProfileToRefresh(pctx, cfg.Grpc.PlayerUrl, &playerPb.FindOnePlayerProfileToRefreshReq{
		PlayerId: strings.TrimPrefix(claims.PlayerId, "player:"),
	})
	if err != nil {
		return nil, err
	}

	accessToken := jwtAuth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtAuth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtAuth.ReloadToken(cfg.Jwt.RefreshSecretKey, claims.ExpiresAt.Unix(), &jwtAuth.Claims{
		PlayerId: profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	if err := u.authRepository.UpdateOnePlayerCredential(pctx, req.CredentialId, &auth.UpdateRefreshTokenReq{
		PlayerId:     profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	credential, err := u.authRepository.FindOnePlayerCredential(pctx, req.CredentialId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileInterceptor{
		PlayerProfile: &player.PlayerProfile{
			Id:        "player:" + profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			Id:           credential.Id.Hex(),
			PlayerId:     credential.PlayerId,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUsecase) Logout(pctx context.Context, credentialId string) (int64, error) {
	return u.authRepository.DeleteOnePlayerCredential(pctx, credentialId)
}

func (u *authUsecase) AccessTokenSearch(pctx context.Context, accessToken string) (*authPb.AccessTokenSearchRes, error) {
	credential, err := u.authRepository.FindOneAccessToken(pctx, accessToken)
	if err != nil {
		return &authPb.AccessTokenSearchRes{
			IsValid: false,
		}, err
	}

	if credential == nil {
		return &authPb.AccessTokenSearchRes{
			IsValid: false,
		}, errors.New("error: access token not found")
	}

	return &authPb.AccessTokenSearchRes{
		IsValid: true,
	}, nil
}

func (u *authUsecase) RoleCount(pctx context.Context) (*authPb.RolesCountRes, error) {
	result, err := u.authRepository.RolesCount(pctx)
	if err != nil {
		return nil, err
	}

	return &authPb.RolesCountRes{
		Count: result,
	}, nil
}
