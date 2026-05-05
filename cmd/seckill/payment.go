package main

import (
	"fmt"
	"seckill/internal/paymentSvr"

	"github.com/spf13/cobra"
)

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "Start PaymentSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting PaymentSvr...")
		paymentSvr.RunPaymentSvr()
	},
}
