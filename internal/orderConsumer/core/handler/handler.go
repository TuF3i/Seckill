package handler

import (
	"context"
	"seckill/internal/orderConsumer/core/models"
	"seckill/internal/orderConsumer/core/pkg/lconst"
	"time"

	"github.com/bytedance/sonic"
)

func (r *Handler) checkTimeout() {
	orders, err := r.D.GetUnpaidOrders()
	if err != nil {
		r.L.Warn("Error fetching unpaid orders: %v", err.Error())
		return
	}

	now := time.Now()
	for _, order := range orders {
		if now.Sub(order.CreatedAt) > lconst.OrderTimeout {
			err = r.D.UpdateOrderStatus(order.OrderId, 3)
			if err != nil {
				r.L.Warn("Error cancelling order %s: %v", order.OrderId, err)
				continue
			}

			err = r.C.SetOrderStatus(context.Background(), order.OrderId, 3)
			if err != nil {
				r.L.Warn("Error updating cache for cancelled order %s: %v", order.OrderId, err)
				continue
			}

			r.L.Info("Order %s auto-cancelled due to timeout", order.OrderId)
		}
	}
}

func (r *Handler) RunTimeoutChecker() {
	go func() {
		r.Wg.Add(1)
		ticker := time.NewTicker(lconst.TimeoutCheckTick)
		defer func() {
			ticker.Stop()
			r.Wg.Done()
		}()

		for {
			select {
			case <-r.Done:
				return
			case <-ticker.C:
				r.checkTimeout()
			}
		}
	}()
}

func (r *Handler) StopTimeoutChecker() {
	r.Done <- struct{}{}
}

func (r *Handler) RunOrderConsumer() {
	r.Wg.Add(1)
	defer r.Wg.Done()
	for {
		select {
		case <-r.Quit:
			return
		default:
		}

		msg, err := r.Consumer.ReadMessage(r.Ctx)
		if err != nil {
			r.L.Warn("Error reading message: %v", err)
			continue
		}

		var orderMsg models.OrderMessage
		err = sonic.Unmarshal(msg.Value, &orderMsg)
		if err != nil {
			r.L.Warn("Error unmarshalling message: %v", err)
			continue
		}

		order := &models.Order{
			OrderId: orderMsg.OrderId,
			UserId:  orderMsg.UserId,
			ItemId:  orderMsg.ItemId,
			Price:   orderMsg.Price,
			Status:  1,
		}

		err = r.D.SaveOrder(order)
		if err != nil {
			r.L.Warn("Error saving order %s: %v", orderMsg.OrderId, err)
			continue
		}

		err = r.C.SetOrderStatus(r.Ctx, orderMsg.OrderId, 1)
		if err != nil {
			r.L.Warn("Error setting order status in cache for %s: %v", orderMsg.OrderId, err)
			continue
		}

		r.L.Info("Order %s saved successfully", orderMsg.OrderId)
	}
}

func (r *Handler) StopOrderConsumer() {
	r.Cancel()
	r.Quit <- struct{}{}
}

func (r *Handler) WaitUntilDone() {
	r.Wg.Wait()
	close(r.Done)
	close(r.Quit)
}
