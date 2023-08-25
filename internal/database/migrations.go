package database

import (
	"reflect"
	"sort"
	"strings"

	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"
)

// IMPORTANT: We do not support removal of fields in models.!

var (
	migrations []*migrate.Migration = []*migrate.Migration{
		{
			ID: hashModel(databasemodels.User{}),
			Migrate: func(e *xorm.Engine) error {
				return e.Sync(&databasemodels.User{})
			},
			Rollback: func(e *xorm.Engine) error {
				return e.DropTables(&databasemodels.User{})
			},
		},
		{
			ID: hashModel(databasemodels.Url{}),
			Migrate: func(e *xorm.Engine) error {
				return e.Sync(&databasemodels.Url{})
			},
			Rollback: func(e *xorm.Engine) error {
				return e.DropTables(&databasemodels.Url{})
			},
		},
	}
)

func hashModel(s any) string {
	v := reflect.ValueOf(s)
	t := v.Type()
	var fieldList []string
	for i := 0; i < v.NumField(); i++ {
		n := t.Field(i).Name
		f := n
		f += "_" + t.Field(i).Type.Name()
		xv, _ := t.Field(i).Tag.Lookup("xorm")
		if xv != "" {
			f += "_" + strings.ReplaceAll(xv, " ", "-")
		}
		fieldList = append(fieldList, f)
	}
	sort.Sort(sort.StringSlice(fieldList))
	return strings.Join(fieldList, "__")
}

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
