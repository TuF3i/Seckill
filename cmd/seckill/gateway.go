package main

import (
	"fmt"
	gateway "seckill/internal/gateway/core/app"

	"github.com/spf13/cobra"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Start API Gateway service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting API Gateway...")
		gateway.OnCreate()
		defer gateway.OnDestory()
		gateway.RunServer()
	},
}
