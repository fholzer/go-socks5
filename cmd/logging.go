package main

import (
	stdlog "log"
	"os"

	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func setupLogging() {
	log = &logrus.Logger{
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		Out: ansicolor.NewAnsiColorWriter(os.Stdout),

		// Log as JSON instead of the default ASCII formatter.
		Formatter: &logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		},
		Hooks: make(logrus.LevelHooks),

		// Only log the warning severity or above.
		Level: logrus.InfoLevel,
	}

	w := logrus.New().Writer()
	//defer w.Close()
	stdlog.SetOutput(w)
}
