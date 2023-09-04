package logger

import (
	"fmt"
	"log"
	"os"
	"plato/common/config"
	"time"
)

var logger log.Logger

func init() {
	logFile, _ := os.OpenFile(
		"./debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644)
	logger.SetOutput(logFile)
}

func Debug(a ...any) {
	if config.GetGlobalEnv() != "debug" {
		return
	}

	prefix := fmt.Sprintf("[%v] ", time.Now())
	logger.Println(prefix, a)
}
