package url

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/constrains"
	"github.com/mhkarimi1383/url-shortener/internal/controller"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"
	requestschemas "github.com/mhkarimi1383/url-shortener/types/request_schemas"
	responseschemas "github.com/mhkarimi1383/url-shortener/types/response_schemas"
)

func Redirect(c echo.Context) error {
	u := databasemodels.Url{
		ShortCode: c.Param(constrains.ShortCodeParamName),
	}
	if has, _ := database.Engine.Get(&u); !has {
		return echo.ErrNotFound
	}
	now := time.Now()
	u.VisitCount++
	u.LastVisitedAt = &now
	u.Entity.VisitCount++
	u.Entity.LastVisitedAt = &now
	database.Engine.ID(u.Id).Update(&u)
	database.Engine.ID(u.Entity.Id).Update(&u.Entity)
	target := u.FullUrl
	if configuration.CurrentConfig.AddRefererQueryParam {
		url, _ := url.Parse(target)
		q := url.Query()
		q.Add(constrains.RefererQueryParam, c.Scheme()+"://"+c.Request().Host+c.Request().URL.Path)
		url.RawQuery = q.Encode()
		target = url.String()
	}
	return c.Redirect(http.StatusTemporaryRedirect, target)
}

func Create(c echo.Context) error {
	user := c.Get(constrains.UserInfoContextVar).(databasemodels.User)

	r := new(requestschemas.CreateURL)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	shortCode, err := controller.CreateUrl(r, user)
	if err != nil {
		return err
	}
	shortURL, err := url.JoinPath(c.Scheme()+"://"+c.Request().Host, configuration.CurrentConfig.BaseURI, "/"+shortCode)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, responseschemas.Create{
		ShortCode: shortCode,
		ShortUrl:  shortURL,
	})
}

func Delete(c echo.Context) error {
	user := c.Get(constrains.UserInfoContextVar).(databasemodels.User)

	id, err := strconv.ParseInt((c.Param(constrains.IdParamName)), 10, 0)
	if err != nil {
		return err
	}
	if err := controller.DeleteUrl(id, user); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func RemoveOldIds(c echo.Context) error {
	user := c.Get(constrains.UserInfoContextVar).(databasemodels.User)

	// Define cutoff date (6 months ago)
	cutoff := time.Now().AddDate(0, -1, 0)

	// Prepare query conditions
	session := database.Engine.Where("last_visited_at < ?", cutoff)

	if !user.Admin {
		// Non-admins can only delete their own URLs
		session = session.And("creator_id = ?", user.Id)
	}

	// Delete matching records
	_, err := session.Delete(&databasemodels.Url{})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

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

	list, err := controller.ListUrls(user, limit, offset)
	if err != nil {
		return err
	}

	for i, item := range list.Result {
		shortURL, err := url.JoinPath(c.Scheme()+"://"+c.Request().Host, configuration.CurrentConfig.BaseURI, "/"+item.ShortCode)
		if err != nil {
			return err
		}
		list.Result[i].ShortUrl = shortURL
	}

	return c.JSON(http.StatusOK, list)
}
