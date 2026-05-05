package main

import (
	"fmt"
	"seckill/internal/itemSvr"

	"github.com/spf13/cobra"
)

var itemCmd = &cobra.Command{
	Use:   "item",
	Short: "Start ItemSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting ItemSvr...")
		itemSvr.RunItemSvr()
	},
}
