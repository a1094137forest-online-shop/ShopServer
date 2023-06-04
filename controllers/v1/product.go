package v1

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes"
	"github.com/lithammer/shortuuid"

	"ShopServer/constant"
	"ShopServer/postgresql"
	"ShopServer/postgresql/model"
	"ShopServer/proto/ShopServer"
)

func (s *ShopServe) GetProduct(ctx context.Context, req *ShopServer.GetProductReq) (*ShopServer.GetProductResp, error) {
	var resp ShopServer.GetProductResp

	productID := req.ProductID

	log.Println("productID", productID)

	p, err := model.GetProduct(ctx, postgresql.PoolWr.Read(), productID)
	if err != nil {
		return nil, err
	}

	timestampProto, _ := ptypes.TimestampProto(*p.Created_at)

	resp = ShopServer.GetProductResp{
		Code: constant.SUCCESS,
		Msg:  "Ok",
		Data: &ShopServer.GetProductRespInfo{
			ProductID:   p.ProductID,
			Title:       p.Title,
			Description: p.Description,
			CreatedAt:   timestampProto,
		},
	}

	return &resp, nil
}

func (s *ShopServe) CreateProduct(ctx context.Context, req *ShopServer.CreateProductReq) (*ShopServer.CreateProductResp, error) {

	shopID, err := model.GetShopIDByAccount(ctx, postgresql.PoolWr.Read(), req.Account)
	if err != nil {
		return nil, err
	}
	var p = model.Product{
		ProductID:   constant.PRODUCT_PREFIX + shortuuid.New(),
		ShopID:      *shopID,
		Title:       req.Data.Title,
		Description: req.Data.Description,
	}

	if err := p.Upsert(ctx, postgresql.PoolWr.Write()); err != nil {
		return nil, err
	}

	var resp = ShopServer.CreateProductResp{
		Code: constant.SUCCESS,
		Msg:  "Ok",
	}
	return &resp, nil
}

func (s *ShopServe) GetProductsList(ctx context.Context, req *ShopServer.GetProductsListReq) (*ShopServer.GetProductsListResp, error) {
	var resp = ShopServer.GetProductsListResp{
		Code: constant.SUCCESS,
		Msg:  "Ok",
	}

	shopID, err := model.GetShopIDByAccount(ctx, postgresql.PoolWr.Read(), req.Account)
	if err != nil {
		return nil, err
	}
	ps, err := model.GetProductsList(ctx, postgresql.PoolWr.Read(), *shopID)
	if err != nil {
		return nil, err
	}

	products := make([]*ShopServer.ProductsListInfo, len(*ps))
	for i, p := range *ps {
		timestampProto, _ := ptypes.TimestampProto(*p.Created_at)
		products[i] = &ShopServer.ProductsListInfo{
			ProductID:   p.ProductID,
			Title:       p.Title,
			Description: p.Description,
			CreatedAt:   timestampProto,
		}
	}

	resp.Data = products
	log.Println("resp.Data", resp.Data)
	return &resp, nil
}

func (s *ShopServe) UpdateProduct(ctx context.Context, req *ShopServer.UpdateProductReq) (*ShopServer.UpdateProductResp, error) {
	var resp ShopServer.UpdateProductResp

	shopID, err := model.GetShopIDByAccount(ctx, postgresql.PoolWr.Read(), req.Account)
	if err != nil {
		return nil, err
	}

	var p = model.Product{
		ProductID:   req.Data.ProductID,
		ShopID:      *shopID,
		Title:       req.Data.Title,
		Description: req.Data.Description,
	}

	if err := p.Upsert(ctx, postgresql.PoolWr.Write()); err != nil {
		return nil, err
	}

	resp = ShopServer.UpdateProductResp{
		Code: constant.SUCCESS,
		Msg:  "Ok",
	}
	return &resp, nil
}

func (s *ShopServe) DeleteProduct(ctx context.Context, req *ShopServer.DeleteProductReq) (*ShopServer.DeleteProductResp, error) {
	var resp ShopServer.DeleteProductResp

	shopID, err := model.GetShopIDByAccount(ctx, postgresql.PoolWr.Read(), req.Account)
	if err != nil {
		return nil, err
	}

	var p = model.Product{
		ProductID: req.Data.ProductID,
		ShopID:    *shopID,
	}

	if err := p.Delete(ctx, postgresql.PoolWr.Write()); err != nil {
		return nil, err
	}

	resp = ShopServer.DeleteProductResp{
		Code: constant.SUCCESS,
		Msg:  "Ok",
	}
	return &resp, nil
}
