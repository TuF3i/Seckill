package main

import (
	"fmt"
	itemSvr "seckill/internal/itemSvr/core/app"

	"github.com/spf13/cobra"
)

var itemCmd = &cobra.Command{
	Use:   "item",
	Short: "Start ItemSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting ItemSvr...")
		itemSvr.OnCreate()
		defer itemSvr.OnDestory()
		itemSvr.RunServer()
	},
}
