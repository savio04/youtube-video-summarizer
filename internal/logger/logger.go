package logger

import "go.uber.org/zap"

var AppLogger *zap.SugaredLogger

func InitAppLogger() {
	AppLogger = zap.Must(zap.NewProduction()).Sugar()
}
