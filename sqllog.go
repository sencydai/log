package log

import (
	"github.com/go-xorm/core"
)

type mysqlLogger struct {
	showSql bool
	level   core.LogLevel
}

var SqlLogger = &mysqlLogger{showSql: true, level: core.LOG_ERR}

func (ml *mysqlLogger) Debug(v ...interface{}) {
	if ml.level >= core.LOG_DEBUG {
		writeBuffer(DEBUG_N, 4, v...)
	}
}

func (ml *mysqlLogger) Debugf(format string, v ...interface{}) {
	if core.LOG_DEBUG >= ml.level {
		writeBufferf(DEBUG_N, 4, format, v...)
	}
}

func (ml *mysqlLogger) Error(v ...interface{}) {
	if core.LOG_ERR >= ml.level {
		writeBuffer(ERROR_N, 4, v...)
	}
}
func (ml *mysqlLogger) Errorf(format string, v ...interface{}) {
	if core.LOG_ERR >= ml.level {
		writeBufferf(ERROR_N, 4, format, v...)
	}
}
func (ml *mysqlLogger) Info(v ...interface{}) {
	if core.LOG_INFO >= ml.level {
		writeBuffer(INFO_N, 4, v...)
	}
}
func (ml *mysqlLogger) Infof(format string, v ...interface{}) {
	if core.LOG_INFO >= ml.level {
		writeBufferf(INFO_N, 4, format, v...)
	}
}
func (ml *mysqlLogger) Warn(v ...interface{}) {
	if core.LOG_WARNING >= ml.level {
		writeBuffer(WARN_N, 4, v...)
	}
}
func (ml *mysqlLogger) Warnf(format string, v ...interface{}) {
	if core.LOG_WARNING >= ml.level {
		writeBufferf(WARN_N, 4, format, v...)
	}
}

func (ml *mysqlLogger) Level() core.LogLevel {
	return ml.level
}

func (ml *mysqlLogger) SetLevel(l core.LogLevel) {
	ml.level = l
}

func (ml *mysqlLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		ml.showSql = show[0]
	} else {
		ml.showSql = true
	}
}

func (ml *mysqlLogger) IsShowSQL() bool {
	return ml.showSql
}
