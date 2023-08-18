package database

import (
	"time"

	"go.uber.org/zap"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"github.com/mhkarimi1383/url-shortener/internal/log"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
)

var Engine *xorm.Engine

func Init() {
	var err error
	Engine, err = xorm.NewEngine(configuration.GetConfig().DatabaseEngine, configuration.GetConfig().DatabaseConnectionString)
	if err != nil {
		log.Logger.Panic(
			err.Error(),
			zap.String("driver", configuration.GetConfig().DatabaseEngine),
			zap.String("connection-string", configuration.GetConfig().DatabaseConnectionString),
		)
	}

	Engine.SetMaxIdleConns(configuration.GetConfig().DatabaseMaxIdleConnections)
	Engine.SetMaxOpenConns(configuration.GetConfig().DatabaseMaxOpenConnections)
	Engine.SetConnMaxLifetime(configuration.GetConfig().DatabaseMaxConnectionLifetime * time.Second)
	Engine.SetLogger(newZapLogger(log.Logger.Sugar()))
}
