package main

import (
	"seckill/internal/initdb"

	"github.com/spf13/cobra"
)

var initDBCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Initialize database tables",
	Long:  "Connect to PostgreSQL and auto-migrate all required database tables (user_table, item_table, order_table)",
	Run: func(cmd *cobra.Command, args []string) {
		initdb.InitDB()
	},
}

func init() {
	rootCmd.AddCommand(initDBCmd)
}
