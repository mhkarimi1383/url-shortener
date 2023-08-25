package controller

import (
	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
	"github.com/mhkarimi1383/url-shortener/types/response_schemas"
	"github.com/mhkarimi1383/url-shortener/utils/shortcode"
)

func CreateURL(r *requestschemas.CreateURL, creator databasemodels.User) (string, error) {
	u := databasemodels.Url{
		FullUrl: r.FullUrl,
		Creator: creator,
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

func ListURLs(userID int64, limit, offset int) (responseschemas.ListURLs, error) {
	user := databasemodels.User{
		Id: userID,
	}
	if _, err := database.Engine.Get(&user); err != nil {
		return nil, err
	}
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
