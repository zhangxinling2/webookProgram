package startup

import (
	"webookProgram/webook/pkg/logger"
)

func InitLogger() logger.LoggerV1 {
	return &logger.NoOpLogger{}
}
