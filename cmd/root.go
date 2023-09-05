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
*/package cmd

import (
	"os"
	"strconv"
	"time"

	"github.com/brpaz/echozap"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/mhkarimi1383/url-shortener/constrains"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/internal/endpoint/entity"
	"github.com/mhkarimi1383/url-shortener/internal/endpoint/url"
	"github.com/mhkarimi1383/url-shortener/internal/endpoint/user"
	"github.com/mhkarimi1383/url-shortener/internal/flagutil"
	"github.com/mhkarimi1383/url-shortener/internal/log"
	ivalidator "github.com/mhkarimi1383/url-shortener/internal/validator"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/ui"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "url-shortener",
		Short: "Simple and minimalism URL Shortener",
		Long:  ``,
		Run:   start,
	}
	cfg configuration.Config
)

func Execute() {
	if invalid := flagutil.SetFlagsFromEnv(rootCmd.PersistentFlags(), "USH"); invalid.String != "" {
		log.Logger.Panic("Invalid environment values provided", invalid)
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
		DurationVar(&cfg.DatabaseMaxConnectionLifetime, "database-max-connection-lifetime", 300*time.Second, "Maximum lifetime for database connections in second used by connection pool")
	rootCmd.PersistentFlags().StringVar(&cfg.DatabaseEngine, "database-engine", "sqlite", "The engine of database")
	rootCmd.PersistentFlags().StringVar(&cfg.DatabaseConnectionString, "database-connection-string", "./database.sqlite3", "Connection string of database")
	rootCmd.PersistentFlags().StringVar(&cfg.JWTSecret, "jwt-secret", "superdupersecret", "jwt secret to sign tokens with, strongly recommended to change")
	rootCmd.PersistentFlags().BoolVar(&cfg.AddRefererQueryParam, "add-referer-query-param", true, "Add 'referer' query param to redirect url or not")
	rootCmd.PersistentFlags().IntVar(&cfg.RandomGeneratorMax, "random-generator-max", 10000, "Generator will use this to generate shortcodes (higher value = bigger shortcodes), at least 10000 should be set")
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

	authMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.RegisteredClaims)
		},
		ContextKey: constrains.UserTokenContextVar,
		SigningKey: []byte(configuration.CurrentConfig.JWTSecret),
	})

	checkUserExists := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Get(constrains.UserTokenContextVar).(*jwt.Token)
			strID := token.Claims.(*jwt.RegisteredClaims).ID
			id, err := strconv.ParseInt(strID, 10, 0)
			if err != nil {
				return err
			}
			usr := databasemodels.User{Id: id}
			if has, _ := database.Engine.Get(&usr); !has {
				return echo.ErrForbidden
			}
			c.Set(constrains.UserInfoContextVar, usr)
			return next(c)
		}
	}

	checkUserAdmin := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			usr := c.Get(constrains.UserInfoContextVar).(databasemodels.User)
			if !usr.Admin {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}

	e.Use(echozap.ZapLogger(log.Logger))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Validator = ivalidator.EchoValidator
	e.HidePort = true
	e.HideBanner = true

	e.Any("/:"+constrains.ShortCodeParamName, url.Redirect)
	e.Any("/:"+constrains.ShortCodeParamName+"/", url.Redirect)

	apiGroup := e.Group("/api")

	userGroup := apiGroup.Group("/user")
	userGroup.POST("/login/", user.Login)
	userGroup.POST("/register/", user.Register)
	userGroup.PUT("/change-password/:"+constrains.IdParamName+"/", user.ChangePassword, authMiddleware, checkUserExists, checkUserAdmin)
	userGroup.POST("/", user.Create, authMiddleware, checkUserExists, checkUserAdmin)

	urlGroup := apiGroup.Group("/url", authMiddleware, checkUserExists)
	urlGroup.POST("/", url.Create)
	urlGroup.GET("/", url.List)
	urlGroup.DELETE("/:"+constrains.IdParamName+"/", url.Delete)

	entityGroup := apiGroup.Group("/entity", authMiddleware, checkUserExists, checkUserAdmin)
	entityGroup.GET("/", entity.List)
	entityGroup.POST("/", entity.Create)
	entityGroup.DELETE("/:"+constrains.IdParamName+"/", entity.Delete)

	uiGroup := e.Group("/ui")
	uiGroup.StaticFS("/", ui.MainFS)
	uiGroup.StaticFS("/assets/", ui.AssetsFS)
	uiGroup.FileFS("/*.html", "index.html", ui.MainFS)
	e.FileFS("/favicon.ico", "favicon.ico", ui.MainFS)
	e.FileFS("/ui/favicon.ico", "favicon.ico", ui.MainFS)

	if configuration.CurrentConfig.RunServer {
		log.Logger.Info("WebServer Started", zap.String("listen-address", configuration.CurrentConfig.ListenAddress))
		log.Logger.Fatal(
			e.Start(configuration.CurrentConfig.ListenAddress).Error(),
			zap.String("listen-address", configuration.CurrentConfig.ListenAddress),
		)
	}
}
