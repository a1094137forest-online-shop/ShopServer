package v1

import (
	"context"
	"log"

	"github.com/lithammer/shortuuid"

	"ShopServer/constant"
	"ShopServer/postgresql"
	"ShopServer/postgresql/model"
	"ShopServer/proto/ShopServer"
)

func (s *ShopServe) CreateShop(ctx context.Context, req *ShopServer.CreateShopReq) (*ShopServer.CreateShopResp, error) {
	var resp = ShopServer.CreateShopResp{
		Code: constant.SUCCESS,
		Msg:  "Ok",
	}

	var shop = model.Shop{
		ShopID:  constant.SHOP_PREFIX + shortuuid.New(),
		Account: req.Account,
	}

	err := shop.Upsert(ctx, postgresql.PoolWr.Write())
	if err != nil {
		log.Println("upser err", err)
		return nil, err
	}

	return &resp, nil
}
