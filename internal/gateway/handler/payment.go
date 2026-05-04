package handler

import (
	"context"
	"seckill/internal/gateway/dto"
	"seckill/internal/gateway/dto/payment"
	"seckill/internal/gateway/pkg/lerror"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func (r *Handler) ProcessPaymentHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req payment.ProcessPaymentReq
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		success, err := r.PaymentSvr.ProcessPayment(ctx, req.OrderId, req.UserId)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		data := payment.ProcessPaymentResp{Success: success}
		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}
