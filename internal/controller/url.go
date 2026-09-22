package controller

import (
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"
	requestschemas "github.com/mhkarimi1383/url-shortener/types/request_schemas"
	responseschemas "github.com/mhkarimi1383/url-shortener/types/response_schemas"
	"github.com/mhkarimi1383/url-shortener/utils/shortcode"
)

func CreateEntity(r *requestschemas.CreateEntity, creator databasemodels.User) error {
	e := databasemodels.Entity{
		Name:        r.Name,
		Description: r.Description,
		Creator:     creator,
	}

	if _, err := database.Engine.Insert(&e); err != nil {
		return err
	}

	return nil
}

func ValidateCreateUrl(r *requestschemas.CreateURL) error {
	if configuration.CurrentConfig.RejectRedirectUrls {
		parsed, _ := url.Parse(r.FullUrl) // URL Already validated
		if !(slices.Contains(configuration.CurrentConfig.WhiteListHosts, parsed.Host)) {
			isRedirect, err := shortcode.IsRedirectingURL(r.FullUrl)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Error while checking URL redirection: "+err.Error())
			}
			if isRedirect {
				return echo.NewHTTPError(http.StatusBadRequest, "Shortened URLs are not allowed.")
			}
		}
	}
	return nil
}

func ValidateUpdateUrl(r *requestschemas.UpdateURL) error {
	if configuration.CurrentConfig.RejectRedirectUrls {
		parsed, _ := url.Parse(r.FullUrl) // URL Already validated
		if !(slices.Contains(configuration.CurrentConfig.WhiteListHosts, parsed.Host)) {
			isRedirect, err := shortcode.IsRedirectingURL(r.FullUrl)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Error while checking URL redirection: "+err.Error())
			}
			if isRedirect {
				return echo.NewHTTPError(http.StatusBadRequest, "Shortened URLs are not allowed.")
			}
		}
	}
	return nil
}

func CreateUrl(r *requestschemas.CreateURL, creator databasemodels.User) (databasemodels.Url, error) {
	entity := databasemodels.Entity{}
	if r.Entity != 0 {
		entity.Id = r.Entity
		has, err := database.Engine.Get(&entity)
		if err != nil {
			return databasemodels.Url{}, err
		}
		if !has {
			return databasemodels.Url{}, echo.NewHTTPError(http.StatusBadRequest, "Entity does not exist.")
		}
	}
	u := databasemodels.Url{
		FullUrl: r.FullUrl,
		Creator: creator,
		Entity:  entity,
	}
	if len(r.ShortCode) > 0 {
		u.ShortCode = r.ShortCode
	} else {
		u.ShortCode = shortcode.Generate(u.Id, time.Now())
	}
	if _, err := database.Engine.Insert(&u); err != nil {
		return databasemodels.Url{}, err
	}
	return u, nil
}

func UpdateUrl(id int64, r *requestschemas.UpdateURL, user databasemodels.User) error {
	condition := databasemodels.Url{Id: id}
	if !user.Admin {
		condition.Creator = user
	}

	affected, err := database.Engine.Table(new(databasemodels.Url)).NoVersionCheck().Incr("version").Cols("full_url").Update(
		&databasemodels.Url{FullUrl: r.FullUrl},
		&condition,
	)
	if err != nil {
		return err
	}
	if affected == 0 {
		return echo.ErrNotFound
	}
	return nil
}

func DeleteUrl(id int64, user databasemodels.User) error {
	u := databasemodels.Url{Id: id}
	if !user.Admin {
		u.Creator = user
	}
	has, err := database.Engine.Get(&u)
	if err != nil {
		return err
	}
	if !has {
		return echo.ErrNotFound
	}
	if _, err := database.Engine.ID(u.Id).Delete(&databasemodels.Url{}); err != nil {
		return err
	}
	return nil
}

func GetUrl(id int64, user databasemodels.User) (*responseschemas.Url, error) {
	u := databasemodels.Url{Id: id}
	if !user.Admin {
		u.Creator = user
	}
	has, err := database.Engine.Get(&u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, echo.ErrNotFound
	}
	return &responseschemas.Url{Url: u}, nil
}

func ListUrls(user databasemodels.User, limit, offset int) (*responseschemas.ListUrls, error) {
	var urls []databasemodels.Url
	prepared := new(responseschemas.ListUrls)
	u := new(databasemodels.Url)
	if !user.Admin {
		u.Creator = user
	}
	if err := database.Engine.Limit(limit, offset).Find(&urls, u); err != nil {
		return nil, err
	}
	total, err := database.Engine.Count(u)
	if err != nil {
		return nil, err
	}
	totalVisit, err := database.Engine.SumInt(u, "visit_count")
	if err != nil {
		return nil, err
	}
	prepared.MetaData.Count = total
	prepared.MetaData.TotalVisit = totalVisit
	for _, u := range urls {
		prepared.Result = append(prepared.Result, responseschemas.Url{Url: u})
	}
	return prepared, nil
}

func ListEntities(user databasemodels.User, limit, offset int) (*responseschemas.ListEntities, error) {
	var entities []databasemodels.Entity
	e := new(databasemodels.Entity)
	if !user.Admin {
		e.Creator = user
	}
	prepared := new(responseschemas.ListEntities)
	if err := database.Engine.Limit(limit, offset).Find(&entities, e); err != nil {
		return nil, err
	}
	prepared.Result = entities
	total, err := database.Engine.Count(e)
	if err != nil {
		return nil, err
	}
	totalVisit, err := database.Engine.SumInt(e, "visit_count")
	if err != nil {
		return nil, err
	}
	prepared.MetaData.Count = total
	prepared.MetaData.TotalVisit = totalVisit
	return prepared, nil
}

func DeleteEntity(id int64) error {
	session := database.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	entity := databasemodels.Entity{
		Id: id,
	}
	has, err := session.Get(&entity)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	if !has {
		_ = session.Rollback()
		return echo.ErrNotFound
	}

	if _, err := session.Delete(&databasemodels.Url{
		Entity: entity,
	}); err != nil {
		_ = session.Rollback()
		return err
	}

	if _, err := session.Delete(&databasemodels.Entity{
		Id: id,
	}); err != nil {
		_ = session.Rollback()
		return err
	}
	if err := session.Commit(); err != nil {
		return err
	}
	return nil
}
