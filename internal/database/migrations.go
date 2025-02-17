package database

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"
	"xorm.io/xorm"
	"xorm.io/xorm/migrate"

	databasemodels "github.com/mhkarimi1383/url-shortener/types/database_models"
)

// IMPORTANT: We do not support removal of fields in models.!

var migrations []*migrate.Migration = []*migrate.Migration{
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
	{
		ID: hashModel(databasemodels.Entity{}),
		Migrate: func(e *xorm.Engine) error {
			return e.Sync(&databasemodels.Entity{})
		},
		Rollback: func(e *xorm.Engine) error {
			return e.DropTables(&databasemodels.Entity{})
		},
	},
}

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
	slices.Sort(fieldList)
	algorithm := md5.New()
	algorithm.Write([]byte(t.Name() + ":" + strings.Join(fieldList, "__")))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func RunMigrations() error {
	m := migrate.New(
		Engine,
		&migrate.Options{
			TableName:    "migration",
			IDColumnName: "id",
		},
		migrations,
	)

	return m.Migrate()
}
