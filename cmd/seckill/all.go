package main

import (
	"fmt"
	"seckill/internal/gateway"
	"seckill/internal/itemSvr"
	"seckill/internal/orderConsumer"
	"seckill/internal/orderSvr"
	"seckill/internal/paymentSvr"
	"seckill/internal/userSvr"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Start all microservices",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting all microservices...")

		go itemSvr.RunItemSvr()
		go orderSvr.RunOrderSvr()
		go orderConsumer.RunOrderConsumer()
		go paymentSvr.RunPaymentSvr()
		go userSvr.RunUserSvr()

		gateway.RunGateway()
	},
}
