package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/internal/controller"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
	"github.com/mhkarimi1383/url-shortener/types/response_schemas"
)

const (
	firstUserAlreadyExist = "First User already exist"
	IdParamName           = "Id"
	UserInfoContextVar    = "userInfo"
	UserTokenContextVar   = "user"
)

func Login(c echo.Context) error {
	l := new(requestschemas.Login)
	if err := c.Bind(l); err != nil {
		return err
	}
	if err := c.Validate(l); err != nil {
		return err
	}

	user, token, err := controller.Login(l)
	if err != nil {
		if errors.Is(err, controller.ErrInvalidUsernameOrPassword) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return err
	}
	return c.JSON(http.StatusOK, responseschemas.Login{
		Token: token,
		Info:  *user,
	})
}

func Register(c echo.Context) error {
	r := new(requestschemas.Register)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	usr := new(databasemodels.User)
	total, err := database.Engine.Count(usr)
	if err != nil {
		return err
	}
	if total > 0 {
		return echo.NewHTTPError(http.StatusConflict, firstUserAlreadyExist)
	}
	if err := controller.CreateUser(r, true); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func Create(c echo.Context) error {
	r := new(requestschemas.CreateUser)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	if err := controller.CreateUser(
		&requestschemas.Register{
			Username: r.Username,
			Password: r.Password,
		},
		r.Admin,
	); err != nil {
		if errors.Is(err, controller.ErrUserAlreadyExist) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func ChangePassword(c echo.Context) error {
	r := new(requestschemas.ChangeUserPassword)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	idStr := c.Param(IdParamName)
	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := controller.ChangeUserPassword(r, id); err != nil {
		if errors.Is(err, controller.ErrUserDoesNotExist) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
