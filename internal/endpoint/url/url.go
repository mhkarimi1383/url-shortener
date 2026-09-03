package url

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/mhkarimi1383/url-shortener/constrains"
	"github.com/mhkarimi1383/url-shortener/internal/controller"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/internal/redirectcache"
	"github.com/mhkarimi1383/url-shortener/internal/visits"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"
	requestschemas "github.com/mhkarimi1383/url-shortener/types/request_schemas"
	responseschemas "github.com/mhkarimi1383/url-shortener/types/response_schemas"
)

func Redirect(c echo.Context) error {
	shortCode := c.Param(constrains.ShortCodeParamName)
	entry, found, err := redirectcache.Default.Resolve(c.Request().Context(), shortCode, func() (redirectcache.Entry, bool, error) {
		u := databasemodels.Url{ShortCode: shortCode}
		has, err := database.Engine.Get(&u)
		if err != nil {
			return redirectcache.Entry{}, false, err
		}
		if !has {
			return redirectcache.Entry{}, false, nil
		}
		return redirectcache.Entry{
			URLID:    u.Id,
			EntityID: u.Entity.Id,
			Target:   u.FullUrl,
		}, true, nil
	})
	if err != nil {
		return err
	}
	if !found {
		return echo.ErrNotFound
	}
	visits.Default.Increment(entry.URLID, entry.EntityID, time.Now())
	return redirect(c, entry.Target)
}

func redirect(c echo.Context, target string) error {
	if configuration.CurrentConfig.AddRefererQueryParam {
		parsed, err := url.Parse(target)
		if err != nil {
			return err
		}
		q := parsed.Query()
		q.Add(constrains.RefererQueryParam, c.Scheme()+"://"+c.Request().Host+c.Request().URL.Path)
		parsed.RawQuery = q.Encode()
		target = parsed.String()
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
	if err := controller.ValidateCreateUrl(r); err != nil {
		return err
	}

	var shortCode string
	err := redirectcache.Default.Mutate(c.Request().Context(), func() error {
		var createErr error
		shortCode, createErr = controller.CreateUrl(r, user)
		return createErr
	})
	if err != nil {
		return cacheMutationError(err)
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
	err = redirectcache.Default.Mutate(c.Request().Context(), func() error {
		return controller.DeleteUrl(id, user)
	})
	if err != nil {
		return cacheMutationError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

func RemoveUnusedUrls(c echo.Context) error {
	user := c.Get(constrains.UserInfoContextVar).(databasemodels.User)
	cutoffStr := c.QueryParam(constrains.CutoffQueryParamName) // Example: 2006-01-02T15:04:05Z07:00
	if cutoffStr == "" {
		cutoffStr = time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	}

	cutoff, err := time.Parse(time.RFC3339, cutoffStr)
	if err != nil {
		return err
	}

	err = redirectcache.Default.Mutate(c.Request().Context(), func() error {
		return visits.Default.SynchronizeCleanup(func() error {
			session := database.Engine.Where("last_visited_at < ?", cutoff)
			if !user.Admin {
				session = session.And("creator_id = ?", user.Id)
			}

			var expired []databasemodels.Url
			if err := session.Cols("id").Find(&expired); err != nil {
				return err
			}
			if len(expired) == 0 {
				return nil
			}
			ids := make([]int64, 0, len(expired))
			for _, item := range expired {
				ids = append(ids, item.Id)
			}
			_, err := database.Engine.In("id", ids).Delete(&databasemodels.Url{})
			return err
		})
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Unable to safely remove old links.").SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func cacheMutationError(err error) error {
	if errors.Is(err, redirectcache.ErrMutationUnavailable) {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Cache coordination is temporarily unavailable.").SetInternal(err)
	}
	return err
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
