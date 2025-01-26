package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	gormloger "gorm.io/gorm/logger"
)

type LogxLogger struct {
	gormloger.Config
}

func NewLogger(cfg gormloger.Config) gormloger.Interface {
	l := &LogxLogger{
		Config: cfg,
	}
	return l
}

// LogMode 设置日志级别
func (l *LogxLogger) LogMode(level gormloger.LogLevel) gormloger.Interface {
	l.LogLevel = level
	return l
}

// Info 记录信息级别的日志
func (l *LogxLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Infof(msg, data...)
}

// Warn 记录警告级别的日志
func (l *LogxLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).WithFields(logx.LogField{Key: "warn", Value: true}).Infof(msg, data...)
}

// Error 记录错误级别的日志
func (l *LogxLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data...)
}

// Trace 记录执行时间和跟踪信息
func (l *LogxLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormloger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormloger.Error && (!errors.Is(err, gormloger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		var content string
		if rows == -1 {
			content = fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			content = fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		logx.WithContext(ctx).Error(content)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormloger.Warn:
		sql, rows := fc()
		var content string
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			content = fmt.Sprintf("[%s] [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			content = fmt.Sprintf("[%s] [%.3fms] [rows:%v] %s", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		logx.WithContext(ctx).Info(content)
	case l.LogLevel == gormloger.Info:
		sql, rows := fc()
		var content string
		if rows == -1 {
			content = fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			content = fmt.Sprintf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		logx.WithContext(ctx).Info(content)
	}
}
