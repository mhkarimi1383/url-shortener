package responseschemas

import "github.com/mhkarimi1383/url-shortener/types/database_models"

type Create struct {
	ShortURL  string
	ShortCode string
}

type ListURLs []databasemodels.Url
