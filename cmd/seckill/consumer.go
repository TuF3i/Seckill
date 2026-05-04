package main

import (
	"fmt"
	orderConsumer "seckill/internal/orderConsumer/core/app"

	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start OrderConsumer service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting OrderConsumer...")
		orderConsumer.OnCreate()
		defer orderConsumer.OnDestory()
		orderConsumer.RunServer()
	},
}
