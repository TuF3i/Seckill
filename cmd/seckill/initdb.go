package main

import (
	"fmt"
	n "seckill/infrastructures/nacos"
	"seckill/internal/initdb"
	"seckill/pkg/config"

	"github.com/spf13/cobra"
)

var (
	nacosHost     string
	nacosPort     uint64
	nacosUser     string
	nacosPassword string
)

var initDBCmd = &cobra.Command{
	Use:   "initdb",
	Short: "Initialize database tables",
	Long:  "Connect to PostgreSQL and auto-migrate all required database tables (user_table, item_table, order_table)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing database...")

		nacosClient, err := n.NewNacosClient(
			n.WithHost(nacosHost),
			n.WithPort(nacosPort),
			n.WithUserName(nacosUser),
			n.WithPassword(nacosPassword),
		)
		if err != nil {
			panic(fmt.Errorf("initdb: create nacos client: %w", err))
		}

		loader, err := config.NewLoader(nacosClient, "seckill", "REDROCK")
		if err != nil {
			panic(fmt.Errorf("initdb: load config: %w", err))
		}

		cfg := loader.GetConfig()
		if err := initdb.InitDB(&cfg.PostgreSQL); err != nil {
			panic(fmt.Errorf("initdb: %w", err))
		}

		fmt.Println("Database initialized successfully")
	},
}

func init() {
	rootCmd.AddCommand(initDBCmd)

	initDBCmd.Flags().StringVar(&nacosHost, "nacos-host", "localhost", "Nacos server host")
	initDBCmd.Flags().Uint64Var(&nacosPort, "nacos-port", 8848, "Nacos server port")
	initDBCmd.Flags().StringVar(&nacosUser, "nacos-user", "nacos", "Nacos username")
	initDBCmd.Flags().StringVar(&nacosPassword, "nacos-password", "nacos", "Nacos password")
}
