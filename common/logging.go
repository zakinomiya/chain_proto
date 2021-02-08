package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/comail/colog"
	"github.com/google/uuid"
)

func LogginSettings(logFile string, logLevel string) {
	l, err := os.OpenFile(filepath.Join("./", os.Getenv("LOG_DIR"), fmt.Sprintf("chain_%.5s.log", uuid.New().String())), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file. err=", err)
	}
	log.SetOutput(io.MultiWriter(l, os.Stdout))

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
