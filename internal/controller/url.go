package controller

import (
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
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
	if _, err := database.Engine.Insert(&u); err != nil {
		return "", err
	}
	shortant := shortcode.Generate(u.Id, u.CreatedAt)
	u.ShortCode = shortant
	if _, err := database.Engine.Update(&u); err != nil {
		return "", err
	}
	return shortant, nil
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

func ListUrls(user databasemodels.User, limit, offset int) ([]databasemodels.Url, error) {
	var urls []databasemodels.Url
	if user.Admin {
		if err := database.Engine.Limit(limit, offset).Find(&urls); err != nil {
			return nil, err
		}
		return urls, nil
	}
	if err := database.Engine.Limit(limit, offset).Find(&urls, &databasemodels.Url{Creator: user}); err != nil {
		return nil, err
	}
	return urls, nil
}

func ListEntities(limit, offset int) ([]databasemodels.Entity, error) {
	var entities []databasemodels.Entity
	if err := database.Engine.Limit(limit, offset).Find(&entities); err != nil {
		return nil, err
	}
	return entities, nil
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
