package item

import "github.com/Applessr/hello-sekai-shop-tutorial/modules/models"

type (
	CreateItemReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
		Damage   int     `json:"damage" validate:"required"`
	}

	ItemShowCase struct {
		ItemId string  `json:"item_id"`
		Title  string  `json:"title"`
		Price  float64 `json:"price"`
		Damage int     `json:"damage"`
		ImgUrl string  `json:"img_url"`
	}

	ItemSearchReq struct {
		Title string `json:"title"`
		models.PaginateReq
	}

	ItemUpdateReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
		Damage   int     `json:"damage" validate:"required"`
	}

	EnableOrDisableItemReq struct {
		UsageStatus bool `json:"status"`
	}
)
