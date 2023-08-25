package responseschemas

import "github.com/mhkarimi1383/url-shortener/types/database_models"

type Login struct {
	Token string
	Info  databasemodels.User
}
