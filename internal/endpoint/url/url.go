package url

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/internal/controller"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
	"github.com/mhkarimi1383/url-shortener/types/response_schemas"
)

const (
	limitQueryParamName  = "limit"
	offsetQueryParamName = "offset"
)

func Redirect(c echo.Context) error {
	u := databasemodels.Url{
		Shortand: c.Param("shortcode"),
	}
	database.Engine.Get(&u)
	return c.Redirect(http.StatusTemporaryRedirect, u.FullUrl)
}

func Create(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	strID := userToken.Claims.(*jwt.RegisteredClaims).ID
	id, err := strconv.ParseInt(strID, 10, 0)
	if err != nil {
		return err
	}

	r := new(requestschemas.CreateURL)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	user := databasemodels.User{
		Id: id,
	}
	if _, err := database.Engine.Get(&user); err != nil {
		return err
	}
	shortcode, err := controller.CreateURL(r, user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, responseschemas.Create{
		ShortCode: shortcode,
		ShortURL:  c.Scheme() + "://" + c.Request().Host + "/" + shortcode,
	})
}

func List(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	strID := userToken.Claims.(*jwt.RegisteredClaims).ID
	id, err := strconv.ParseInt(strID, 10, 0)
	if err != nil {
		return err
	}

	limitStr := c.QueryParam(limitQueryParamName)
	if limitStr == "" {
		limitStr = "10"
	}
	offsetStr := c.QueryParam(offsetQueryParamName)
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

	list, err := controller.ListURLs(id, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}
