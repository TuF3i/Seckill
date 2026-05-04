package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seckill",
	Short: "Seckill system unified entry point",
	Long:  "A unified command-line entry point for all microservices in the seckill system",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(itemCmd)
	rootCmd.AddCommand(orderCmd)
	rootCmd.AddCommand(paymentCmd)
	rootCmd.AddCommand(consumerCmd)
	rootCmd.AddCommand(gatewayCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(allCmd)
}
