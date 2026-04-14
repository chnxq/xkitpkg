package asynq

import (
	"fmt"

	"github.com/chnxq/xkitmod/log"
	"github.com/hibiken/asynq"
)

const (
	logKey = "[" + KindAsynq + "]"
)

type logger struct {
}

func newLogger() asynq.Logger {
	return &logger{}
}

func (l logger) Debug(args ...any) {
	log.Debugf("%s %s", logKey, fmt.Sprint(args...))
}

func (l logger) Info(args ...any) {
	log.Infof("%s %s", logKey, fmt.Sprint(args...))
}

func (l logger) Warn(args ...any) {
	log.Warnf("%s %s", logKey, fmt.Sprint(args...))
}

func (l logger) Error(args ...any) {
	log.Errorf("%s %s", logKey, fmt.Sprint(args...))
}

func (l logger) Fatal(args ...any) {
	log.Fatalf("%s %s", logKey, fmt.Sprint(args...))
}
