package itemUsecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemRepository"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/models"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	ItemUsecaseService interface {
		CreatedItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		FindOneItem(pctx context.Context, itemId string) (*item.ItemShowCase, error)
		FindManyItem(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error)
		EditItem(pctx context.Context, itemId string, req *item.ItemUpdateReq) (*item.ItemShowCase, error)
		EnableOrDisableItem(pctx context.Context, itemId string) (bool, error)
	}

	itemUsecase struct {
		itemRepository itemRepository.ItemRepositoryService
	}
)

func NewItemUsecase(itemRepository itemRepository.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{itemRepository}
}

func (u *itemUsecase) CreatedItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error) {
	if !u.itemRepository.IsUniqueItem(pctx, req.Title) {
		return nil, errors.New("error: item already exists")
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	itemId, err := u.itemRepository.InsertOneItem(pctx, &item.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		UsageStatus: true,
		ImageUrl:    req.ImageUrl,
		CreatedAt:   utils.LocalTime().In(loc),
		UpdatedAt:   utils.LocalTime().In(loc),
	})
	if err != nil {
		return nil, errors.New("error: insert item failed")
	}

	return u.FindOneItem(pctx, itemId.Hex())
}

func (u *itemUsecase) FindOneItem(pctx context.Context, itemId string) (*item.ItemShowCase, error) {
	result, err := u.itemRepository.FindOneItem(pctx, itemId)
	if err != nil {
		return nil, errors.New("error: find one item not found")
	}
	return &item.ItemShowCase{
		ItemId:   result.Id.Hex(),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}

func (u *itemUsecase) FindManyItem(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error) {
	findItemsFilter := bson.D{}
	findItemOptions := make([]*options.FindOptions, 0)

	countItemsFilter := bson.D{}
	// filter
	if req.Start != "" && req.Limit != 0 {
		req.Start = strings.TrimPrefix(req.Start, "item:")
		findItemsFilter = append(findItemsFilter, bson.E{"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}
	if req.Title != "" {
		findItemsFilter = append(findItemsFilter, bson.E{"title", primitive.Regex{Pattern: req.Title, Options: "i"}})
		countItemsFilter = append(countItemsFilter, bson.E{"title", primitive.Regex{Pattern: req.Title, Options: "i"}})
	}

	findItemsFilter = append(findItemsFilter, bson.E{"usage_status", true})
	countItemsFilter = append(countItemsFilter, bson.E{"usage_status", true})

	//Option
	findItemOptions = append(findItemOptions, options.Find().SetSort(bson.D{{"_id", 1}}))
	findItemOptions = append(findItemOptions, options.Find().SetLimit(int64(req.Limit)))

	//Find
	result, err := u.itemRepository.FindManyItems(pctx, findItemsFilter, findItemOptions)
	if err != nil {
		return nil, errors.New("error: find many items failed")
	}

	total, err := u.itemRepository.CountItems(pctx, countItemsFilter)
	if err != nil {
		return nil, errors.New("error: count items failed")
	}

	if len(result) == 0 {
		return &models.PaginateRes{
			Data:  make([]*item.ItemShowCase, 0),
			Total: total,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil

	}

	return &models.PaginateRes{
		Data:  result,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
		},
		Next: models.NextPaginate{
			Start: result[len(result)-1].ItemId,
			Href:  fmt.Sprintf("%s?limit=%d&title=%s&start=%s", basePaginateUrl, req.Limit, req.Title, result[len(result)-1].ItemId),
		},
	}, nil
}

func (u *itemUsecase) EditItem(pctx context.Context, itemId string, req *item.ItemUpdateReq) (*item.ItemShowCase, error) {
	// Update logical
	updateReq := bson.M{}
	if req.Title != "" {
		if !u.itemRepository.IsUniqueItem(pctx, req.Title) {
			log.Println("Error: EditItem failed: this title is already exist")
			return nil, errors.New("error: this title is already exist")
		}

		updateReq["title"] = req.Title
	}
	if req.ImageUrl != "" {
		updateReq["image_url"] = req.ImageUrl
	}
	if req.Damage > 0 {
		updateReq["damage"] = req.Damage
	}
	if req.Price >= 0 {
		updateReq["price"] = req.Price
	}
	updateReq["updated_at"] = utils.LocalTime()

	if err := u.itemRepository.UpdateOneItem(pctx, itemId, updateReq); err != nil {
		return nil, err
	}

	return u.FindOneItem(pctx, itemId)
}

func (u *itemUsecase) EnableOrDisableItem(pctx context.Context, itemId string) (bool, error) {
	result, err := u.itemRepository.FindOneItem(pctx, itemId)
	if err != nil {
		return false, err
	}

	if err := u.itemRepository.UpdateOneItem(pctx, itemId, bson.M{"usage_status": !result.UsageStatus}); err != nil {
		return false, err
	}

	return !result.UsageStatus, nil
}
