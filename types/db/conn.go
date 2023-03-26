/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package db

import (
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"

	"github.com/mhkarimi1383/url-shortener/config"
)

var engine *xorm.Engine

func GetEngine() *xorm.Engine {
	if engine == nil {
		var err error
		cfg := config.GetConfig()
		engine, err = xorm.NewEngine(cfg.DBEngine, cfg.DBConnectionString)
		if err != nil {
			panic(err)
		}
		engine.ShowSQL(cfg.Debug)
		engine.SetMapper(names.SameMapper{})
		engine.SetLogLevel(log.LOG_DEBUG)

		engine.Sync(new(URL), new(User))
	}
	return engine
}
