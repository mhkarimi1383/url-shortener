/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/types/db"
	"github.com/mhkarimi1383/url-shortener/types/request"
	"github.com/mhkarimi1383/url-shortener/utils/errors"
	"github.com/mhkarimi1383/url-shortener/utils/url"
)

func ListURLs(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return errors.ErrorToHTTPErrorAdapter(http.StatusBadRequest, "limit should be valid number", err)
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return errors.ErrorToHTTPErrorAdapter(http.StatusBadRequest, "offset should be valid number", err)
	}

	urls, err := db.ListURLs(limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &urls)
}

func ListUserURLs(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return errors.ErrorToHTTPErrorAdapter(http.StatusBadRequest, "limit should be valid number", err)
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return errors.ErrorToHTTPErrorAdapter(http.StatusBadRequest, "offset should be valid number", err)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*request.JWTClaims)
	urls, err := db.ListURLsByUserID(limit, offset, claims.UserID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &urls)
}

func CreateURL(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*request.JWTClaims)

	input := request.URL{}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if !url.IsValidUrl(input.UpstreamURL) {
		return errors.ErrorToHTTPErrorAdapter(http.StatusBadRequest, "UpstreamURL is not a valid URL", fmt.Errorf("upstreamURL was is not a valid url"))
	}

	u := &db.URL{
		Creator:       claims.UserID,
		UpstreamURL:   input.UpstreamURL,
		DownStreamURI: fmt.Sprintf("%x", sha256.Sum256([]byte(input.UpstreamURL)))[:8],
	}
	if err := u.Insert(); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &u)
}

func RemoveURL(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.ErrorToHTTPErrorAdapter(http.StatusBadRequest, "id should be valid number", err)
	}

	u := &db.URL{
		ID: id,
	}
	if err := u.Remove(); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
