package entity

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/constrains"
	"github.com/mhkarimi1383/url-shortener/internal/controller"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
)

func List(c echo.Context) error {
	user := c.Get(constrains.UserInfoContextVar).(databasemodels.User)

	limitStr := c.QueryParam(constrains.LimitQueryParamName)
	if limitStr == "" {
		limitStr = "10"
	}
	offsetStr := c.QueryParam(constrains.OffsetQueryParamName)
	if offsetStr == "" {
		offsetStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	list, err := controller.ListEntities(user, limit, offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, list)
}

func Create(c echo.Context) error {
	user := c.Get(constrains.UserInfoContextVar).(databasemodels.User)

	r := new(requestschemas.CreateEntity)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	if err := controller.CreateEntity(r, user); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func Delete(c echo.Context) error {
	idStr := c.Param(constrains.IdParamName)
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := controller.DeleteEntity(id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
