package main

import (
	"fmt"
	orderSvr "seckill/internal/orderSvr/core/app"

	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Start OrderSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting OrderSvr...")
		orderSvr.OnCreate()
		defer orderSvr.OnDestory()
		orderSvr.RunServer()
	},
}
