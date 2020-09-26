package utils

import (
	"log"
	"os"

	"github.com/comail/colog"
)

func LogginSettings(logFile string, logLevel string) {
	level := logLevel

	// Override the minimum log level if LOG_LEVEL is presented
	if value, isPresented := os.LookupEnv("LOG_LEVEL"); isPresented {
		level = value
	}

	if level, err := colog.ParseLevel(level); err != nil {
		colog.SetMinLevel(colog.LInfo)
	} else {
		colog.SetMinLevel(level)
	}

	colog.SetDefaultLevel(colog.LInfo)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
}
