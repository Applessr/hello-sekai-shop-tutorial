package playerUsecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/modules/player"
	playerPb "github.com/Applessr/hello-sekai-shop-tutorial/modules/player/playerPb"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/player/playerRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error)
		FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfile, error)
		AddPlayerMoney(pctx context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error)
		GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredential(pctx context.Context, password, email string) (*playerPb.PlayerProfile, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, playerId string) (*playerPb.PlayerProfile, error)
	}

	playerUsecase struct {
		playerRepository playerRepository.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepository playerRepository.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{playerRepository}
}

func (u *playerUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.playerRepository.GetOffset(pctx)
}
func (u *playerUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.playerRepository.UpserOffset(pctx, offset)
}

func (u *playerUsecase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error) {
	if !u.playerRepository.IsUniquePlayer(pctx, req.Email, req.Username) {
		return nil, errors.New("error: email or username already exists")
	}

	//HashPassword
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error: hash password failed")
	}

	//insert player
	playerId, err := u.playerRepository.InsertOnePlayer(pctx, &player.Player{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Username:  req.Username,
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		PlayerRole: []player.PlayerRole{
			{
				RoleTitle: "player",
				RoleCode:  0,
			},
		},
	})
	if err != nil {
		return nil, errors.New("error: insert player failed")
	}

	return u.FindOnePlayerProfile(pctx, playerId.Hex())
}

func (u *playerUsecase) FindOnePlayerProfile(pctx context.Context, playerId string) (*player.PlayerProfile, error) {
	result, err := u.playerRepository.FindOnePlayerProfile(pctx, playerId)
	if err != nil {
		return nil, errors.New("error: find one player profile not found")
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &player.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.In(loc),
		UpdatedAt: result.UpdatedAt.In(loc),
	}, nil
}

func (u *playerUsecase) AddPlayerMoney(pctx context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error) {
	if err := u.playerRepository.InsertOnePlayerTransaction(pctx, &player.PlayerTransaction{
		PlayerId:  req.PlayerId,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	}); err != nil {
		return nil, errors.New("error: insert player transaction failed")
	}

	return u.playerRepository.GetPlayerSavingAccount(pctx, req.PlayerId)
}

func (u *playerUsecase) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	return u.playerRepository.GetPlayerSavingAccount(pctx, playerId)
}

func (u *playerUsecase) FindOnePlayerCredential(pctx context.Context, password, email string) (*playerPb.PlayerProfile, error) {
	result, err := u.playerRepository.FindOnePlayerCredential(pctx, email)
	if err != nil {
		return nil, errors.New("error: find one player email not found")
	}
	if result == nil {
		return nil, errors.New("player not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		log.Printf("Error: FindOnePlayerCredential: %s", err.Error())
		return nil, errors.New("error: password is invalid")
	}

	roleCode := 0
	for _, v := range result.PlayerRole {
		roleCode += v.RoleCode
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &playerPb.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}

func (u *playerUsecase) FindOnePlayerProfileToRefresh(pctx context.Context, playerId string) (*playerPb.PlayerProfile, error) {
	result, err := u.playerRepository.FindOnePlayerProfileToRefresh(pctx, playerId)
	if err != nil {
		return nil, errors.New("error: find one player profile not found")
	}

	roleCode := 0
	for _, v := range result.PlayerRole {
		roleCode += v.RoleCode
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &playerPb.PlayerProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}
