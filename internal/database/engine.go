package database

import (
	"time"

	"go.uber.org/zap"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"github.com/mhkarimi1383/url-shortener/internal/log"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
)

var Engine *xorm.Engine

func Init() {
	var err error
	Engine, err = xorm.NewEngine(configuration.CurrentConfig.DatabaseEngine, configuration.CurrentConfig.DatabaseConnectionString)
	if err != nil {
		log.Logger.Panic(
			err.Error(),
			zap.String("driver", configuration.CurrentConfig.DatabaseEngine),
			zap.String("connection-string", configuration.CurrentConfig.DatabaseConnectionString),
		)
	}

	Engine.SetMapper(names.GonicMapper{})
	Engine.SetMaxIdleConns(configuration.CurrentConfig.DatabaseMaxIdleConnections)
	Engine.SetMaxOpenConns(configuration.CurrentConfig.DatabaseMaxOpenConnections)
	Engine.SetConnMaxLifetime(configuration.CurrentConfig.DatabaseMaxConnectionLifetime * time.Second)
	Engine.SetLogger(newZapLogger(log.Logger.Sugar()))
}
