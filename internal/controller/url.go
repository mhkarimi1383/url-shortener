package controller

import (
	"time"

	"github.com/mhkarimi1383/url-shortener/internal/database"
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

func CreateUrl(r *requestschemas.CreateURL, creator databasemodels.User) (string, error) {
	entity := databasemodels.Entity{}
	if r.Entity != 0 {
		entity.Id = r.Entity
		_, err := database.Engine.Get(&entity)
		if err != nil {
			return "", err
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
		println(u.ShortCode)
	}
	if _, err := database.Engine.Insert(&u); err != nil {
		return "", err
	}
	return u.ShortCode, nil
}

func DeleteUrl(id int64, user databasemodels.User) error {
	if user.Admin {
		_, err := database.Engine.Delete(&databasemodels.Url{
			Id: id,
		})
		return err
	}
	_, err := database.Engine.Delete(&databasemodels.Url{
		Id:      id,
		Creator: user,
	})
	return err
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
	entity := databasemodels.Entity{
		Id: id,
	}
	if _, err := database.Engine.Get(&entity); err != nil {
		return err
	}

	if _, err := database.Engine.Delete(&databasemodels.Url{
		Entity: entity,
	}); err != nil {
		return err
	}

	if _, err := database.Engine.Delete(&databasemodels.Entity{
		Id: id,
	}); err != nil {
		return err
	}
	return nil
}
