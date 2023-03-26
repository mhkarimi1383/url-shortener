/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package api

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mhkarimi1383/url-shortener/api/handlers"
	"github.com/mhkarimi1383/url-shortener/config"
	"github.com/mhkarimi1383/url-shortener/types/db"
	"github.com/mhkarimi1383/url-shortener/types/request"
)

var EchoInstance *echo.Echo

func init() {
	// Echo instance
	EchoInstance = echo.New()
}

func Serve(listenAddress string) {
	// Middleware
	EchoInstance.Use(middleware.Recover())

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(request.JWTClaims)
		},
		SigningKey: []byte(config.GetConfig().JWTSecret),
	}
	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	// Routes
	EchoInstance.GET("/api/healthz/", handlers.Healthz)
	EchoInstance.GET("/api/status/", handlers.Status)
	EchoInstance.POST("/api/register/", handlers.Register)
	EchoInstance.POST("/api/login/", handlers.Login)
	EchoInstance.GET("/api/user/urls/", handlers.ListUserURLs, jwtMiddleware)
	EchoInstance.GET("/api/urls/", handlers.ListURLs, jwtMiddleware)
	EchoInstance.POST("/api/user/urls/", handlers.CreateURL, jwtMiddleware)
	EchoInstance.DELETE("/api/user/urls/:id/", handlers.RemoveURL, jwtMiddleware)

	EchoInstance.Any("/*", func(c echo.Context) error {
		u := &db.URL{
			DownStreamURI: strings.TrimPrefix(strings.TrimPrefix(c.Request().RequestURI, "/"), "/"),
		}
		println(c.Request().RequestURI)
		_, err := u.Get()
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusTemporaryRedirect, u.UpsteamURL)
	})

	// Start server
	EchoInstance.Logger.Fatal(EchoInstance.Start(listenAddress))
}
