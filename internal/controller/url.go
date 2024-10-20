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
	if user.Admin {
		if err := database.Engine.Limit(limit, offset).Find(&urls); err != nil {
			return nil, err
		}
		total, err := database.Engine.Count(new(databasemodels.Url))
		if err != nil {
			return nil, err
		}
		for _, u := range urls {
			prepared.Result = append(prepared.Result, responseschemas.Url{Url: u})
		}
		prepared.MetaData.Count = total
		return prepared, nil
	}
	if err := database.Engine.Limit(limit, offset).Find(&urls, &databasemodels.Url{Creator: user}); err != nil {
		return nil, err
	}
	total, err := database.Engine.Count(&databasemodels.Url{
		Creator: user,
	})
	if err != nil {
		return nil, err
	}
	for _, u := range urls {
		prepared.Result = append(prepared.Result, responseschemas.Url{Url: u})
	}
	prepared.MetaData.Count = total
	return prepared, nil
}

func ListEntities(limit, offset int) (*responseschemas.ListEntities, error) {
	var entities []databasemodels.Entity
	prepared := new(responseschemas.ListEntities)
	if err := database.Engine.Limit(limit, offset).Find(&entities); err != nil {
		return nil, err
	}
	prepared.Result = entities
	total, err := database.Engine.Count(new(databasemodels.Entity))
	if err != nil {
		return nil, err
	}
	prepared.MetaData.Count = total
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
