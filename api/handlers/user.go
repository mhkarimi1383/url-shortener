/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/mhkarimi1383/url-shortener/config"
	"github.com/mhkarimi1383/url-shortener/types/db"
	"github.com/mhkarimi1383/url-shortener/types/request"
	"github.com/mhkarimi1383/url-shortener/utils/errors"
)

func Register(c echo.Context) error {
	count, err := db.CountUsers()
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.ErrorToHTTPErrorAdapter(
			http.StatusForbidden,
			"You are not allowed to register",
			fmt.Errorf("user wants to register but he/she was not the first user"),
		)
	}

	input := request.User{}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if len(input.Name) > 60 {
		return errors.ErrorToHTTPErrorAdapter(
			http.StatusForbidden,
			"Username should be lower than 60 characters",
			fmt.Errorf("entered username was too large"),
		)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		if err == bcrypt.ErrPasswordTooLong {
			return errors.ErrorToHTTPErrorAdapter(
				http.StatusBadRequest,
				"Password too large",
				err,
			)
		} else {
			return err
		}
	}

	user := &db.User{
		Name:     input.Name,
		Password: string(password),
	}

	err = user.Insert()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func Login(c echo.Context) error {
	input := request.User{}

	if err := c.Bind(&input); err != nil {
		return err
	}

	u := &db.User{
		Name: input.Name,
	}
	if _, err := u.Get(); err != nil {
		return err
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.ErrorToHTTPErrorAdapter(
				http.StatusBadRequest,
				"Invalid Credentials",
				err,
			)
		} else {
			return err
		}
	}
	claims := request.JWTClaims{
		Name:   input.Name,
		UserID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(config.GetConfig().TokenExpireTime))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &request.TokenResponse{
		Token:      t,
		ExpireTime: time.Now().Add(time.Hour * time.Duration(config.GetConfig().TokenExpireTime)).UnixMilli(),
		TokenType:  "Bearer",
	})
}
