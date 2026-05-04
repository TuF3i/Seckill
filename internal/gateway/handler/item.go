package handler

import (
	"context"
	"seckill/internal/gateway/dto"
	"seckill/internal/gateway/dto/item"
	"seckill/internal/gateway/pkg/lcontext"
	"seckill/internal/gateway/pkg/lerror"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func (r *Handler) AddItemHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req item.AddItemReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		itemId, err := r.ItemSvr.AddItem(ctx, req.Name, req.Stock, req.Price, req.Description)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		data := item.AddItemResp{ItemId: itemId}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}

func (r *Handler) DeleteItemHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req item.DeleteItemReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		err := r.ItemSvr.DeleteItem(ctx, req.ItemId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, nil))
	}
}

func (r *Handler) ListItemsHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, err := lcontext.GetClaimsFromRequestContext(c)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		items, err := r.ItemSvr.ListItems(ctx, claims.UID, claims.Role.String())
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		var data []item.ItemInfoResp
		for _, i := range items {
			data = append(data, item.ItemInfoResp{
				ID:          i.ID,
				Name:        i.Name,
				Stock:       i.Stock,
				Price:       i.Price,
				Description: i.Description,
			})
		}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}
func (r *Handler) StartFlashSaleHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req item.FlashSaleReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		err := r.ItemSvr.StartFlashSale(ctx, req.ItemId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, nil))
	}
}

func (r *Handler) StopFlashSaleHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req item.FlashSaleReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		err := r.ItemSvr.StopFlashSale(ctx, req.ItemId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, nil))
	}
}
