package main

import (
	"fmt"
	"seckill/internal/orderConsumer"

	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start OrderConsumer service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting OrderConsumer...")
		orderConsumer.RunOrderConsumer()
	},
}
