package logger

import (
	"io"
	"os"

	"github.com/longzhoufeng/go-core/debug/writer"
	log "github.com/longzhoufeng/go-logger"
	"github.com/longzhoufeng/go-logger/zap"
	"github.com/longzhoufeng/go-sdk/pkg"
)

// SetupLogger 日志 cap 单位为kb
func SetupLogger(opts ...Option) log.Logger {
	op := setDefault()
	for _, o := range opts {
		o(&op)
	}
	if !pkg.PathExist(op.path) {
		err := pkg.PathCreate(op.path)
		if err != nil {
			log.Fatalf("create dir error: %s", err.Error())
		}
	}
	var err error
	var output io.Writer
	switch op.stdout {
	case "file":
		output, err = writer.NewFileWriter(
			writer.WithPath(op.path),
			writer.WithCap(op.cap<<10),
		)
		if err != nil {
			log.Fatal("logger setup error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}
	var level log.Level
	level, err = log.GetLevel(op.level)
	if err != nil {
		log.Fatalf("get logger level error, %s", err.Error())
	}

	switch op.driver {
	case "zap":
		log.DefaultLogger, err = zap.NewLogger(log.WithLevel(level), log.WithOutput(output), zap.WithCallerSkip(2))
		if err != nil {
			log.Fatalf("new zap logger error, %s", err.Error())
		}
	//case "logrus":
	//	setLogger = logrus.NewLogger(logger.WithLevel(level), logger.WithOutput(output), logrus.ReportCaller())
	default:
		log.DefaultLogger = log.NewLogger(log.WithLevel(level), log.WithOutput(output))
	}
	return log.DefaultLogger
}
