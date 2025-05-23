package authRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/auth"
	playerPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/player/playerPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/grpccon"
	jwtAuth "github.com/Applessr/hello-sekai-shop-tutorial/pkg/jwtauth"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthRepositoryService interface {
		CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error)
		InsertOneCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error)
		FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error)
		UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error
		DeleteOnePlayerCredential(pctx context.Context, credentialId string) (int64, error)
		FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error)
		RolesCount(pctx context.Context) (int64, error)
		AccessToken(cfg *config.Config, claims *jwtAuth.Claims) string
		RefreshToken(cfg *config.Config, claims *jwtAuth.Claims) string
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{db}
}

func (r *authRepository) authDbConnect(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}

func (r *authRepository) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtAuth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRpc client connection failed: %s", err.Error())
		return nil, errors.New("error: gRpc client connection failed")
	}

	result, err := conn.Player().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: CredentialSearch: %s", err.Error())
		return nil, errors.New(err.Error())
	}

	return result, nil
}

func (r *authRepository) InsertOneCredential(pctx context.Context, req *auth.Credential) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConnect(ctx)
	col := db.Collection("auth")

	result, err := col.InsertOne(ctx, req)
	if err != nil || result.InsertedID == nil {
		log.Printf("Error: InsertOneCredential: %s", err.Error())
		return primitive.NilObjectID, errors.New("error: Insert one credential failed")
	}
	log.Printf("Success: InsertOneCredential: %s", result.InsertedID)

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *authRepository) FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConnect(ctx)
	col := db.Collection("auth")

	result := new(auth.Credential)

	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)}).Decode(result); err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Error: FindOnePlayerCredential: %s", err.Error())
		return nil, errors.New("error: find one player credential not found")
	}

	return result, nil
}

func (r *authRepository) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	jwtAuth.SetApiKeyInContext(&ctx)
	conn, err := grpccon.NewGrpcClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gRpc client connection failed: %s", err.Error())
		return nil, errors.New("error: gRpc client connection failed")
	}

	result, err := conn.Player().FindOnePlayerProfileToRefresh(ctx, req)
	if err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh: %s", err.Error())
		return nil, errors.New(err.Error())
	}

	return result, nil
}

func (r *authRepository) UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConnect(ctx)
	col := db.Collection("auth")

	result, err := col.UpdateOne(
		ctx,
		bson.M{"_id": utils.ConvertToObjectId(credentialId)},
		bson.M{
			"$set": bson.M{
				"player_id":     req.PlayerId,
				"access_token":  req.AccessToken,
				"refresh_token": req.RefreshToken,
				"updated_at":    utils.LocalTime(),
			},
		},
	)
	if err != nil || result.ModifiedCount == 0 {
		log.Printf("Error: UpdateOnePlayerCredential: %s", err.Error())
		return errors.New("error: update one player credential failed")
	}

	return nil
}

func (r *authRepository) DeleteOnePlayerCredential(pctx context.Context, credentialId string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConnect(ctx)
	col := db.Collection("auth")

	result, err := col.DeleteOne(ctx, bson.M{"_id": utils.ConvertToObjectId(credentialId)})
	if err != nil || result.DeletedCount == 0 {
		log.Printf("Error: DeleteOnePlayerCredential: %s", err.Error())
		return -1, errors.New("error: delete one player credential failed")
	}
	log.Printf("Success: DeleteOnePlayerCredential: %s", result.DeletedCount)

	return result.DeletedCount, nil
}

func (r *authRepository) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConnect(ctx)
	col := db.Collection("auth")

	credential := new(auth.Credential)
	if err := col.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(credential); err != nil {
		log.Printf("Error: FindOneAccessToken: %s", err.Error())
		return nil, errors.New("error: find one access token not found")
	}

	return credential, nil
}

func (r *authRepository) RolesCount(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConnect(ctx)
	col := db.Collection("roles")

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil && count == 0 {
		log.Printf("Error: RolesCount: %s", err.Error())
		return -1, errors.New("error: roles count failed")
	}

	return count, nil
}

func (r *authRepository) AccessToken(cfg *config.Config, claims *jwtAuth.Claims) string {
	return jwtAuth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtAuth.Claims{
		PlayerId: claims.PlayerId,
		RoleCode: int(claims.RoleCode),
	}).SignToken()
}

func (r *authRepository) RefreshToken(cfg *config.Config, claims *jwtAuth.Claims) string {
	return jwtAuth.NewRefreshToken(cfg.Jwt.RefreshSecretKey, cfg.Jwt.RefreshDuration, &jwtAuth.Claims{
		PlayerId: claims.PlayerId,
		RoleCode: int(claims.RoleCode),
	}).SignToken()
}
