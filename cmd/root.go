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
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/http"
	neturl "net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	"github.com/mhkarimi1383/url-shortener/internal/redirectcache"
	ivalidator "github.com/mhkarimi1383/url-shortener/internal/validator"
	"github.com/mhkarimi1383/url-shortener/internal/visits"
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
	rootCmd.PersistentFlags().StringVarP(&cfg.RootRedirect, "root-redirect", "r", "/BASE_URI/ui/", "Path/URL to redirect when someone comes to root url")
	rootCmd.PersistentFlags().StringVar(&cfg.BaseURI, "base-uri", "/", "Base URL of the project")
	rootCmd.PersistentFlags().BoolVar(&cfg.RejectRedirectUrls, "reject-redirect-urls", false, "Reject already shortened (redirecting) URLs from being stored")
	rootCmd.PersistentFlags().StringSliceVar(&cfg.WhiteListHosts, "white-list-hosts", nil, "Whitelist for redirecting URLs rejection")
	rootCmd.PersistentFlags().StringVar(&cfg.RedisAddress, "redis-address", "", "Redis host:port used for redirect caching; empty disables Redis")
	rootCmd.PersistentFlags().StringVar(&cfg.RedisUsername, "redis-username", "", "Redis username")
	rootCmd.PersistentFlags().StringVar(&cfg.RedisPassword, "redis-password", "", "Redis password")
	rootCmd.PersistentFlags().IntVar(&cfg.RedisDatabase, "redis-database", 0, "Redis database number")
	rootCmd.PersistentFlags().BoolVar(&cfg.RedisTLS, "redis-tls", false, "Use TLS for Redis")
	rootCmd.PersistentFlags().DurationVar(&cfg.RedisCacheTTL, "redis-cache-ttl", 24*time.Hour, "TTL for cached redirects")
	rootCmd.PersistentFlags().DurationVar(&cfg.RedisDialTimeout, "redis-dial-timeout", 100*time.Millisecond, "Redis dial timeout")
	rootCmd.PersistentFlags().DurationVar(&cfg.RedisReadTimeout, "redis-read-timeout", 100*time.Millisecond, "Redis read timeout")
	rootCmd.PersistentFlags().DurationVar(&cfg.RedisWriteTimeout, "redis-write-timeout", 100*time.Millisecond, "Redis write timeout")
	rootCmd.PersistentFlags().DurationVar(&cfg.VisitFlushInterval, "visit-flush-interval", 5*time.Second, "Interval for flushing aggregated visits")
	rootCmd.PersistentFlags().IntVar(&cfg.VisitBufferMaxEntries, "visit-buffer-max-entries", 100000, "Maximum unique URLs buffered for visit statistics")
}

func start(_ *cobra.Command, _ []string) {
	log.Logger.Info("Setting and Validating configuration parameters")

	cfg.BaseURI = strings.TrimSuffix(cfg.BaseURI, "/")
	cfg.RootRedirect = strings.ReplaceAll(cfg.RootRedirect, "/BASE_URI/", cfg.BaseURI+"/")

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
	redirectcache.Init(configuration.CurrentConfig)
	visits.Init(configuration.CurrentConfig)

	if configuration.CurrentConfig.Migrate {
		log.Logger.Info("Running database migrations")
		if err := database.RunMigrations(); err != nil {
			log.Logger.Fatal("Database migration failed", zap.Error(err))
		}
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

	addTrailingSlashMiddleware := middleware.AddTrailingSlashWithConfig(
		middleware.TrailingSlashConfig{
			RedirectCode: http.StatusPermanentRedirect,
			Skipper:      middleware.DefaultSkipper,
		},
	)

	e.Use(echozap.ZapLogger(log.Logger))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Validator = ivalidator.EchoValidator
	e.HidePort = true
	e.HideBanner = true

	rootGroup := e.Group(cfg.BaseURI)

	rootGroup.Any("/:"+constrains.ShortCodeParamName, url.Redirect)
	rootGroup.Any("/:"+constrains.ShortCodeParamName+"/", url.Redirect)

	rootGroup.Any("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, cfg.RootRedirect)
	})

	apiGroup := rootGroup.Group("/api")
	apiGroup.Use(addTrailingSlashMiddleware)

	userGroup := apiGroup.Group("/user")
	userGroup.POST("/login/", user.Login)
	userGroup.POST("/register/", user.Register)
	userGroup.GET("/", user.List, authMiddleware, checkUserExists, checkUserAdmin)
	userGroup.PUT("/change-password/:"+constrains.IdParamName+"/", user.AdminChangePassword, authMiddleware, checkUserExists, checkUserAdmin)
	userGroup.PUT("/change-password/", user.ChangePassword, authMiddleware, checkUserExists)
	userGroup.POST("/", user.Create, authMiddleware, checkUserExists, checkUserAdmin)
	// TODO: logout API (token revoke)

	urlGroup := apiGroup.Group("/url", authMiddleware, checkUserExists)
	urlGroup.POST("/", url.Create)
	urlGroup.GET("/", url.List)
	urlGroup.DELETE("/:"+constrains.IdParamName+"/", url.Delete)

	entityGroup := apiGroup.Group("/entity", authMiddleware, checkUserExists, checkUserAdmin)
	entityGroup.GET("/", entity.List)
	entityGroup.POST("/", entity.Create)
	entityGroup.DELETE("/:"+constrains.IdParamName+"/", entity.Delete)

	uiGroup := rootGroup.Group("/ui")
	uiGroup.Any("", nil, addTrailingSlashMiddleware)
	uiCacheDir := filepath.Join(os.TempDir(), "url-shortener-ui-cache")
	if err := os.MkdirAll(uiCacheDir, 0755); err != nil {
		log.Logger.Panic("Failed to create UI cache directory", zap.Error(err))
	}
	cacheLock := sync.Mutex{}
	uiGroup.GET("/*", func(c echo.Context) error {
		prefix := filepath.Join("/"+strings.TrimPrefix(cfg.BaseURI, "/"), "/ui/")
		filePath := strings.TrimPrefix(strings.TrimPrefix(c.Request().URL.Path, prefix), "/")
		fileNameMd5 := md5.Sum([]byte(filePath))
		fileNameHash := hex.EncodeToString(fileNameMd5[:])
		cachedFilePath := filepath.Join(uiCacheDir, fileNameHash)
		fileParts := strings.Split(filePath, ".")
		mimeType := mime.TypeByExtension("." + fileParts[len(fileParts)-1])
		if mimeType == "" {
			mimeType = "text/html; charset=utf-8"
		}
		if _, err := os.Stat(cachedFilePath); err == nil {
			stream, err := os.Open(cachedFilePath)
			if err != nil {
				return err
			}
			defer stream.Close()
			return c.Stream(200, mimeType, stream)
		}
		err := func() error {
			cacheLock.Lock()
			defer cacheLock.Unlock()
			file, err := ui.MainFS.Open(filePath)
			if err != nil {
				filePath = "index.html"
				file, _ = ui.MainFS.Open(filePath)
			}
			defer file.Close()
			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, file)
			if err != nil {
				return err
			}
			newURI, err := neturl.JoinPath("/BASE_URI/", "../"+cfg.BaseURI+"/")
			if err != nil {
				return err
			}
			modifiedContent := []byte(strings.ReplaceAll(buf.String(), "/BASE_URI/", newURI))
			if err := os.WriteFile(cachedFilePath, modifiedContent, 0644); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
		stream, err := os.Open(cachedFilePath)
		if err != nil {
			return err
		}
		defer stream.Close()
		return c.Stream(200, mimeType, stream)
	})

	commandGroup := apiGroup.Group("/command", authMiddleware, checkUserExists)
	commandGroup.DELETE("/remove-old-links", url.RemoveUnusedUrls)

	if configuration.CurrentConfig.RunServer {
		if err := runServer(e); err != nil {
			log.Logger.Fatal("WebServer stopped unexpectedly", zap.Error(err))
		}
	}
}

func runServer(e *echo.Echo) error {
	workerContext, stopWorker := context.WithCancel(context.Background())
	workerDone := make(chan struct{})
	go func() {
		defer close(workerDone)
		visits.Default.Run(workerContext)
	}()

	serverError := make(chan error, 1)
	go func() {
		log.Logger.Info("WebServer Started", zap.String("listen-address", configuration.CurrentConfig.ListenAddress))
		serverError <- e.Start(configuration.CurrentConfig.ListenAddress)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	var runError error
	select {
	case err := <-serverError:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			runError = err
		}
	case received := <-signals:
		log.Logger.Info("Shutting down", zap.String("signal", received.String()))
		shutdownContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := e.Shutdown(shutdownContext); err != nil {
			log.Logger.Error("Failed to gracefully shut down web server", zap.Error(err))
		}
		cancel()
	}

	stopWorker()
	<-workerDone
	if err := redirectcache.Default.Close(); err != nil {
		log.Logger.Warn("Failed to close Redis", zap.Error(err))
	}
	return runError
}
