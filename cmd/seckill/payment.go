package main

import (
	"fmt"
	paymentSvr "seckill/internal/paymentSvr/core/app"

	"github.com/spf13/cobra"
)

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "Start PaymentSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting PaymentSvr...")
		paymentSvr.OnCreate()
		defer paymentSvr.OnDestory()
		paymentSvr.RunServer()
	},
}
