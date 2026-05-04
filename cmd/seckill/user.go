package main

import (
	"fmt"
	userSvr "seckill/internal/userSvr/core/app"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Start UserSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting UserSvr...")
		userSvr.OnCreate()
		defer userSvr.OnDestory()
		userSvr.RunServer()
	},
}
