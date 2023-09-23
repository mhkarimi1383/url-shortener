package responseschemas

import databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"

type Create struct {
	ShortUrl  string
	ShortCode string
}

type Url struct {
	databasemodels.Url `json:",inline"`
	ShortUrl           string
}

type ListUrls struct {
	Result   []Url
	MetaData MetaData
}

type ListEntities struct {
	MetaData MetaData
	Result   []databasemodels.Entity
}
