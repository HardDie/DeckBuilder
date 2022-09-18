package logger

import (
	"log"
	"os"
)

var (
	Debug = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime)
	Info  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	Warn  = log.New(os.Stdout, "WARN\t", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)
