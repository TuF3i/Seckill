package main

import (
	"fmt"
	gateway "seckill/internal/gateway"
	itemSvr "seckill/internal/itemSvr/core/app"
	orderConsumer "seckill/internal/orderConsumer/core/app"
	orderSvr "seckill/internal/orderSvr/core/app"
	paymentSvr "seckill/internal/paymentSvr/core/app"
	userSvr "seckill/internal/userSvr/core/app"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Start all microservices",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting all microservices...")

		go func() {
			userSvr.OnCreate()
			defer userSvr.OnDestory()
			userSvr.RunServer()
		}()

		go func() {
			itemSvr.OnCreate()
			defer itemSvr.OnDestory()
			itemSvr.RunServer()
		}()

		go func() {
			orderSvr.OnCreate()
			defer orderSvr.OnDestory()
			orderSvr.RunServer()
		}()

		go func() {
			orderConsumer.OnCreate()
			defer orderConsumer.OnDestory()
			orderConsumer.RunServer()
		}()

		go func() {
			paymentSvr.OnCreate()
			defer paymentSvr.OnDestory()
			paymentSvr.RunServer()
		}()

		gateway.OnCreate()
		defer gateway.OnDestory()
		gateway.RunServer()
	},
}
