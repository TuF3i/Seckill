package main

import (
	"fmt"
	"seckill/internal/userSvr"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Start UserSvr service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting UserSvr...")
		userSvr.RunUserSvr()
	},
}
