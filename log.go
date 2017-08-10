package main

import "github.com/sirupsen/logrus"

type Logger struct {
	*logrus.Logger
}

func (log *Logger) Critical(args ...interface{})                 { log.Error(args...) }
func (log *Logger) Criticalf(format string, args ...interface{}) { log.Errorf(format, args...) }
func (log *Logger) Notice(args ...interface{})                   { log.Info(args...) }
func (log *Logger) Noticef(format string, args ...interface{})   { log.Infof(format, args...) }
