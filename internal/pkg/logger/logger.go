package logger

import (
	"log"
	"os"

	"enterprise-orbit/configs"
)

func Init(cfg *configs.Config) error {
	level := cfg.Log.Level
	log.Printf("Logger initialized with level: %s", level)
	return nil
}

func Info(msg string) {
	log.Println("[INFO]", msg)
}

func Error(msg string) {
	log.Println("[ERROR]", msg)
}

func Fatal(msg string) {
	log.Println("[FATAL]", msg)
	os.Exit(1)
}
