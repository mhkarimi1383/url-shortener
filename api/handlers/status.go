/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/types/db"
	"github.com/mhkarimi1383/url-shortener/utils/errors"
)

func Healthz(c echo.Context) error {
	err := db.GetEngine().Ping()
	if err != nil {
		return errors.ErrorToHTTPErrorAdapter(http.StatusInternalServerError, "Error while pinging database", err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "I'm OK.!",
	})
}

// client uses this api to check if register is required or not
func Status(c echo.Context) error {
	count, err := db.CountUsers()
	if err != nil {
		return err
	}

	urlCount, err := db.CountUsers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"first_run":      count == 0,
		"number_of_urls": urlCount,
	})
}
