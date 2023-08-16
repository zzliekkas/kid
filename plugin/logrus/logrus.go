package logrus

import (
	"context"
	"fmt"

	"github.com/leor-w/kid/config"
	"github.com/leor-w/kid/logger"
	"github.com/leor-w/kid/plugin"
	"github.com/leor-w/kid/plugin/logrus/formatters"
	"github.com/leor-w/kid/plugin/logrus/rotates/simple"
	"github.com/leor-w/kid/utils"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	opts   Options
	logger *logrus.Entry
}

func (log *logrusLogger) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(plugin.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	return Default(confName)
}

type Option func(*Options)

func (log *logrusLogger) Init(opts ...logger.Option) error {
	logrusLog := logrus.NewEntry(logrus.New())
	for _, o := range opts {
		o(&log.opts.Options)
	}

	if formatter, ok := log.opts.Context.Value(formatterKey{}).(logrus.Formatter); ok {
		log.opts.formatter = formatter
	}

	if reportCall, ok := log.opts.Context.Value(reportCallKey{}).(bool); ok {
		log.opts.reportCall = reportCall
	}

	logrusLog.Logger.SetReportCaller(log.opts.reportCall)
	logrusLog.Logger.SetLevel(loggerToLogrusLevel(log.opts.Level))
	logrusLog.Logger.SetFormatter(log.opts.formatter)
	log.logger = logrusLog
	return nil
}

func (log *logrusLogger) Options() *logger.Options {
	return &log.opts.Options
}

func (log *logrusLogger) WithFields(fields map[string]interface{}) {
	log.logger.WithFields(fields)
}

func (log *logrusLogger) Hook(hook logger.Hook) {
	log.logger.Logger.AddHook(hook)
}

func (log *logrusLogger) Log(level logger.Level, args ...interface{}) {
	log.logger.Log(loggerToLogrusLevel(level), args...)
}

func (log *logrusLogger) Logf(level logger.Level, fmt string, args ...interface{}) {
	log.logger.Logf(loggerToLogrusLevel(level), fmt, args...)
}

func loggerToLogrusLevel(level logger.Level) logrus.Level {
	switch level {
	case logger.TraceLevel:
		return logrus.TraceLevel
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	case logger.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func logrusToLoggerLevel(level logrus.Level) logger.Level {
	switch level {
	case logrus.TraceLevel:
		return logger.TraceLevel
	case logrus.DebugLevel:
		return logger.DebugLevel
	case logrus.InfoLevel:
		return logger.InfoLevel
	case logrus.WarnLevel:
		return logger.WarnLevel
	case logrus.ErrorLevel:
		return logger.ErrorLevel
	case logrus.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}

func NewLogger(opts ...logger.Option) *logrusLogger {
	options := Options{
		Options: logger.Options{
			Level:   logger.InfoLevel,
			LogName: config.GetString("logger.name"),
			Context: context.Background(),
		},
		reportCall: false,
	}
	l := &logrusLogger{
		opts: options,
	}
	_ = l.Init(opts...)
	return l
}

func Default(name string) logger.Logger {
	confPrefix := utils.GetConfigurationItem("logger", name)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config not found configuration item [%s]", confPrefix))
	}
	log := NewLogger(
		WithFormatter(formatters.NewSimpleFormatter()),
		WithReportCall(config.GetBool(utils.GetConfigurationItem(confPrefix, "reportCaller"))),
	)
	rotate, _ := simple.NewSimpleRotate(
		simple.WithPath(config.GetString(utils.GetConfigurationItem(confPrefix, "path"))),
		simple.WithLogName(config.GetString(utils.GetConfigurationItem(confPrefix, "name"))),
		simple.WithLink(config.GetString(utils.GetConfigurationItem(confPrefix, "link"))),
		simple.WithRotate(config.GetDuration(utils.GetConfigurationItem(confPrefix, "rotate"))),
		simple.WithMaxAge(config.GetDuration(utils.GetConfigurationItem(confPrefix, "maxAge"))),
	)
	log.Hook(rotate)
	return log
}
