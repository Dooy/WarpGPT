package test

import "WarpGPT/pkg/logger"

func RegRun() error {
	logger.Log.Debug("hello word")
	return nil
}
func Log(opt ...interface{}) {
	logger.Log.Debug(opt...)
}
