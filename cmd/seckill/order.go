package main

import (
	"fmt"
	"seckill/internal/orderSvr"

	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Start OrderSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting OrderSvr...")
		orderSvr.RunOrderSvr()
	},
}
