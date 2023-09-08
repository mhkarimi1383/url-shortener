package database

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"xorm.io/xorm/log"
)

type zapLogger struct {
	logger *zap.SugaredLogger

	level   log.LogLevel
	showSQL bool
}

func newZapLogger(logger *zap.SugaredLogger) *zapLogger {
	var lvl log.LogLevel
	switch logger.Level() {
	case zapcore.DebugLevel:
		lvl = log.LOG_DEBUG
	case zapcore.InfoLevel:
		lvl = log.LOG_INFO
	case zapcore.ErrorLevel:
		lvl = log.LOG_ERR
	default:
		lvl = log.LOG_UNKNOWN
	}

	return &zapLogger{
		logger:  logger,
		level:   lvl,
		showSQL: true,
	}
}

func (zl *zapLogger) Debug(v ...any) {
	if zl.level <= log.LOG_DEBUG {
		zl.logger.Debug(v...)
	}
}

func (zl *zapLogger) Debugf(format string, v ...any) {
	if zl.level <= log.LOG_DEBUG {
		zl.logger.Debugf(format, v...)
	}
}

func (zl *zapLogger) Error(v ...any) {
	if zl.level <= log.LOG_ERR {
		zl.logger.Error(v...)
	}
}

func (zl *zapLogger) Errorf(format string, v ...any) {
	if zl.level <= log.LOG_ERR {
		zl.logger.Errorf(format, v...)
	}
}

func (zl *zapLogger) Info(v ...any) {
	if zl.level <= log.LOG_INFO {
		zl.logger.Info(v...)
	}
}

func (zl *zapLogger) Infof(format string, v ...any) {
	if zl.level <= log.LOG_INFO {
		zl.logger.Infof(format, v...)
	}
}

func (zl *zapLogger) Warn(v ...any) {
	if zl.level <= log.LOG_WARNING {
		zl.logger.Warn(v...)
	}
}

func (zl *zapLogger) Warnf(format string, v ...any) {
	if zl.level <= log.LOG_WARNING {
		zl.logger.Warnf(format, v...)
	}
}

func (zl *zapLogger) Level() log.LogLevel {
	return zl.level
}

func (zl *zapLogger) SetLevel(l log.LogLevel) {
	zl.level = l
}

func (zl *zapLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		zl.showSQL = true
		return
	}
	zl.showSQL = show[0]
}

func (zl *zapLogger) IsShowSQL() bool {
	return zl.showSQL
}
