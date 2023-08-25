package database

import (
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

var (
	migrations []*migrate.Migration = []*migrate.Migration{
			{
				ID: "create_user_table",
				Migrate: func(e *xorm.Engine) error {
					return e.Sync(&databasemodels.User{})
				},
				Rollback: func(e *xorm.Engine) error {
					return e.DropTables(&databasemodels.User{})
				},
			},
			{
				ID: "create_url_table",
				Migrate: func(e *xorm.Engine) error {
					return e.Sync(&databasemodels.Url{})
				},
				Rollback: func(e *xorm.Engine) error {
					return e.DropTables(&databasemodels.Url{})
				},
			},
		}
)

func RunMigrations() error {
	m := migrate.New(
		Engine,
		&migrate.Options{
			TableName:    "Migration",
			IDColumnName: "Id",
		},
		migrations,
	)

	return m.Migrate()
}
