package itemHandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item"
	"github.com/Applessr/hello-sekai-shop-tutorial/modules/item/itemUsecase"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/request"
	"github.com/Applessr/hello-sekai-shop-tutorial/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	ItemHttpHandlerService interface {
		CreatedItem(c echo.Context) error
		FindOneItem(c echo.Context) error
		FindManyItem(c echo.Context) error
		EditItem(c echo.Context) error
		EnableOrDisableItem(c echo.Context) error
	}

	itemHttpHandler struct {
		cfg         *config.Config
		itemUsecase itemUsecase.ItemUsecaseService
	}
)

func NewItemHttpHandler(cfg *config.Config, itemUsecase itemUsecase.ItemUsecaseService) ItemHttpHandlerService {
	return &itemHttpHandler{cfg, itemUsecase}
}

func (h *itemHttpHandler) CreatedItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.CreateItemReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.CreatedItem(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *itemHttpHandler) FindOneItem(c echo.Context) error {
	ctx := context.Background()

	itemId := c.Param("item_id")

	res, err := h.itemUsecase.FindOneItem(ctx, itemId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) FindManyItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemSearchReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.FindManyItem(ctx, h.cfg.Paginate.ItemNextPageBasedUrl, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) EditItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemUpdateReq)
	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.EditItem(ctx, itemId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}
	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) EnableOrDisableItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	res, err := h.itemUsecase.EnableOrDisableItem(ctx, itemId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("itemId: %s, status: %v", itemId, res),
	})
}
