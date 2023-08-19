/*
Copyright © 2023 Muhammed Hussein Karimi info@karimi.dev

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/brpaz/echozap"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/internal/flagutil"
	"github.com/mhkarimi1383/url-shortener/internal/log"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Simple and minimalism URL Shortener",
	Long:  ``,
	Run:   start,
}

var cfg configuration.Config

func Execute() {
	if invalid := flagutil.SetFlagsFromEnv(rootCmd.PersistentFlags(), "USH"); invalid.String != "" {
		log.Logger.Panic("Invalid environemt values provided", invalid)
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Logger.Panic(err.Error())
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg.ListenAddress, "listen-address", "l", "127.0.0.1:8080", "Host:Port to listen")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Migrate, "migrate", "m", true, "To run migrations or not")
	rootCmd.PersistentFlags().BoolVarP(&cfg.RunServer, "run-server", "s", true, "To run webserver or not")
	rootCmd.PersistentFlags().IntVar(&cfg.DatabaseMaxIdleConnections, "database-max-idle-connections", 4, "Number of maximum idle connections to database used by connection pool")
	rootCmd.PersistentFlags().IntVar(&cfg.DatabaseMaxOpenConnections, "database-max-open-connections", 10, "Number of maximum open connections to database used by connection pool")
	rootCmd.PersistentFlags().
		DurationVar(&cfg.DatabaseMaxConnectionLifetime, "database-max-connection-lifetime", 300, "Maximum lifetime for database connections in second used by connection pool")
	rootCmd.PersistentFlags().StringVar(&cfg.DatabaseEngine, "database-engine", "sqlite", "The engine of database")
	rootCmd.PersistentFlags().StringVar(&cfg.DatabaseConnectionString, "database-connection-string", "./database.sqlite3", "Connection string of database")
}

func start(_ *cobra.Command, _ []string) {
	log.Logger.Info("Setting and Validating configuration parameters")
	if err := configuration.SetConfig(&cfg); err != nil {
		if vErrs, ok := err.(validator.ValidationErrors); ok {
			for _, vErr := range vErrs {
				log.Logger.Error(
					"Invalid configuration parameter value",
					zap.String("namespace", vErr.Namespace()),
					zap.String("field", vErr.Field()),
					zap.String("struct-namespace", vErr.StructNamespace()),
					zap.String("struct-field", vErr.StructField()),
					zap.String("tag", vErr.Tag()),
					zap.String("actual-tag", vErr.ActualTag()),
					zap.String("kind", vErr.Kind().String()),
					zap.String("type", vErr.Type().String()),
					zap.Any("value", vErr.Value()),
					zap.String("param", vErr.Param()),
				)
			}
			os.Exit(1)
		}
		log.Logger.Panic(err.Error())
	}

	log.Logger.Info("Initializing database engine")
	database.Init()

	if configuration.CurrentConfig.Migrate {
		log.Logger.Info("Running database migrations")
		database.RunMigrations()
	}

	e := echo.New()

	e.Use(echozap.ZapLogger(log.Logger))
	e.Use(middleware.Recover())
	e.HidePort = true
	e.HideBanner = true

	if configuration.CurrentConfig.RunServer {
		log.Logger.Info("Starting WebServer")
		log.Logger.Fatal(
			e.Start(configuration.CurrentConfig.ListenAddress).Error(),
			zap.String("listen-address", configuration.CurrentConfig.ListenAddress),
		)
	}
}
