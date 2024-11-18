package avebi

import "log"

var pkgLogger Logger = log.Default()

type Logger interface {
	Printf(format string, v ...any)
}

func SetLogger(logger Logger) {
	pkgLogger = logger
}
