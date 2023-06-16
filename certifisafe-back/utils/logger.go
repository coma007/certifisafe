package utils

import (
	"fmt"
	"log"
	"os"
)

const (
	Info    = "INFO"
	Warning = "WARNING"
	Error   = "ERROR"
)

const (
	Green  = "\033[1;32m%s\033[0m"
	Yellow = "\033[1;33m%s\033[0m"
	Red    = "\033[1;31m%s\033[0m"
)

func LogInfo(message string) {
	logToFileAndConsole(message, Info, Green)
}

func LogWarning(message string) {
	logToFileAndConsole(message, Warning, Yellow)
}

func LogError(message string) {
	logToFileAndConsole(message, Error, Red)
}

func logToFileAndConsole(message, level, color string) {
	logMessage := fmt.Sprintf("[%s] %s", level, message)
	log.Println(fmt.Sprintf(color, logMessage))

	file, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(logMessage)
}
