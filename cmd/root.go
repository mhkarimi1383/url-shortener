/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mhkarimi1383/url-shortener/api"
	"github.com/mhkarimi1383/url-shortener/config"
	"github.com/mhkarimi1383/url-shortener/types"
	"github.com/mhkarimi1383/url-shortener/types/db"
	"github.com/mhkarimi1383/url-shortener/utils/flagloader"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Simple/Minimalistic URL Shortener",
	Long:  ``,
	Run:   start,
}

func Execute() {
	if err := flagloader.SetFlagsFromEnv(rootCmd.PersistentFlags(), "US"); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

var cfg types.Config

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.ListenAddress, "listen-address", "127.0.0.1:8080", "Address:Port combination used to serve API")
	rootCmd.PersistentFlags().StringVar(&cfg.DBConnectionString, "db-connection-string", "./shortener.db", "Connection string used to connect to database")
	rootCmd.PersistentFlags().StringVar(&cfg.DBEngine, "db-engine", "sqlite3", "Engine for database should be one of (postgres, mysql, sqlite3, mssql)")
	rootCmd.PersistentFlags().StringVar(&cfg.JWTSecret, "jwt-secret", "superdupersecret", "Secret used to sign JWT tokens")
	rootCmd.PersistentFlags().Int64Var(&cfg.TokenExpireTime, "token-expire-time", 5, "How much a token could be valid (In hour)")
	rootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, "Weather to enable debug log or not (will print sensitive data)")
}

func start(_ *cobra.Command, _ []string) {
	config.SetConfig(&cfg)
	defer db.GetEngine().Close()
	api.Serve(cfg.ListenAddress)
}
