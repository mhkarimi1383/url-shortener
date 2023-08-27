package url

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"net/url"

	"github.com/mhkarimi1383/url-shortener/internal/controller"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
	"github.com/mhkarimi1383/url-shortener/types/response_schemas"
)

const (
	limitQueryParamName  = "Limit"
	offsetQueryParamName = "Offset"
	ShortCodeParamName   = "shortCode"
	IdParamName          = "Id"
	refererQueryParam    = "referer"
)

func Redirect(c echo.Context) error {
	u := databasemodels.Url{
		ShortCode: c.Param(ShortCodeParamName),
	}
	if has, _ := database.Engine.Get(&u); !has {
		return echo.ErrNotFound
	}
	target := u.FullUrl
	if configuration.CurrentConfig.AddRefererQueryParam {
		url, _ := url.Parse(target)
		q := url.Query()
		q.Add(refererQueryParam, c.Scheme()+"://"+c.Request().Host+c.Request().URL.Path)
		url.RawQuery = q.Encode()
		target = url.String()
	}
	return c.Redirect(http.StatusTemporaryRedirect, target)
}

func Create(c echo.Context) error {
	user := c.Get("userInfo").(databasemodels.User)

	r := new(requestschemas.CreateURL)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	shortcode, err := controller.CreateUrl(r, user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, responseschemas.Create{
		ShortCode: shortcode,
		ShortUrl:  c.Scheme() + "://" + c.Request().Host + "/" + shortcode,
	})
}

func Delete(c echo.Context) error {
	user := c.Get("userInfo").(databasemodels.User)

	id, err := strconv.ParseInt((c.Param(IdParamName)), 10, 0)
	if err != nil {
		return err
	}
	return controller.DeleteUrl(id, user)
}

func List(c echo.Context) error {
	user := c.Get("userInfo").(databasemodels.User)

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

	var resp responseschemas.ListUrls

	list, err := controller.ListUrls(user.Id, limit, offset)
	if err != nil {
		return err
	}

	for _, item := range list {
		resp = append(resp, responseschemas.Url{
			Url:      item,
			ShortUrl: c.Scheme() + "://" + c.Request().Host + "/" + item.ShortCode,
		})
	}

	return c.JSON(http.StatusOK, resp)
}
