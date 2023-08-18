package database

import (
  "xorm.io/xorm/migrate"
)

func RunMigrations() error {
	m := migrate.New(
	  Engine,
		&migrate.Options{
		  TableName: "Migrations",
			IDColumnName: "Id",
		},
		nil,
	)

  return m.Migrate()
}
