package responseschemas

import databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"

type Login struct {
	Token string
	Info  databasemodels.User
}

type UserList struct {
	MetaData MetaData
	Result   []databasemodels.User
}
