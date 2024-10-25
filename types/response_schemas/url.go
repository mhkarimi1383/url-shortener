package responseschemas

import databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"

type UrlMetaData struct {
	MetaData   `json:",inline"`
	TotalVisit int64
}

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
	MetaData UrlMetaData
}

type ListEntities struct {
	MetaData UrlMetaData
	Result   []databasemodels.Entity
}
