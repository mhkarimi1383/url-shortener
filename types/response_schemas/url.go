package responseschemas

import "github.com/mhkarimi1383/url-shortener/types/database_models"

type Create struct {
	ShortUrl  string
	ShortCode string
}

type Url struct {
	databasemodels.Url `json:",embed"`
	ShortUrl           string
}

type ListUrls []Url
