package handler

import (
	"context"
	"seckill/internal/gateway/dto"
	"seckill/internal/gateway/dto/order"
	"seckill/internal/gateway/pkg/lcontext"
	"seckill/internal/gateway/pkg/lerror"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func (r *Handler) CreateOrderHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims, err := lcontext.GetClaimsFromRequestContext(c)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		var req order.CreateOrderReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		price, err := r.ItemSvr.PrepareOrder(ctx, claims.UID, req.ItemId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		orderId, err := r.OrderSvr.CreateOrder(ctx, claims.UID, req.ItemId, price)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		data := order.CreateOrderResp{OrderId: orderId}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}

func (r *Handler) QueryPaidOrdersHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req order.QueryOrdersReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		orders, err := r.OrderSvr.QueryPaidOrders(ctx, req.UserId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		var data []order.OrderInfoResp
		for _, o := range orders {
			data = append(data, order.OrderInfoResp{
				OrderId:    o.OrderId,
				UserId:     o.UserId,
				ItemId:     o.ItemId,
				Price:      o.Price,
				Status:     int32(o.Status),
				CreateTime: o.CreateTime,
			})
		}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}

func (r *Handler) QueryUnpaidOrdersHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req order.QueryOrdersReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		orders, err := r.OrderSvr.QueryUnpaidOrders(ctx, req.UserId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		var data []order.OrderInfoResp
		for _, o := range orders {
			data = append(data, order.OrderInfoResp{
				OrderId:    o.OrderId,
				UserId:     o.UserId,
				ItemId:     o.ItemId,
				Price:      o.Price,
				Status:     int32(o.Status),
				CreateTime: o.CreateTime,
			})
		}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}

func (r *Handler) QueryCancelledOrdersHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req order.QueryOrdersReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		orders, err := r.OrderSvr.QueryCancelledOrders(ctx, req.UserId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		var data []order.OrderInfoResp
		for _, o := range orders {
			data = append(data, order.OrderInfoResp{
				OrderId:    o.OrderId,
				UserId:     o.UserId,
				ItemId:     o.ItemId,
				Price:      o.Price,
				Status:     int32(o.Status),
				CreateTime: o.CreateTime,
			})
		}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}
