package zapLog

import (
	myZap "github.com/aaronchen2k/deeptest/pkg/core/zap"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// level 日志级别
var level zapcore.Level

// Init 初始化日志服务
func Init() {
	var logger *zap.Logger

	logDir := "log"

	if !dir.IsExist(logDir) {
		dir.InsureDir(logDir)
	}

	level = zap.DebugLevel

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(myZap.GetEncoderCore(level), zap.AddStacktrace(level))
	} else {
		logger = zap.New(myZap.GetEncoderCore(level))
	}
	logger = logger.WithOptions(zap.AddCaller())

	logUtils.Logger = logger
}
