package utils

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
)

const (
	Success = "SUCCESS"
	Info    = "INFO"
	Error   = "ERROR"
)

const (
	Green  = "\033[1;32m%s\033[0m"
	Yellow = "\033[1;33m%s\033[0m"
	Red    = "\033[1;31m%s\033[0m"
)

func LogSuccess(message string, function string) {
	logToFileAndConsole(message, Success, Green, function)
}

func LogInfo(message string, function string) {
	logToFileAndConsole(message, Info, Yellow, function)
}

func LogError(message string, function string) {
	logToFileAndConsole(message, Error, Red, function)
}

func logToFileAndConsole(message, level, color, function string) {
	logMessage := fmt.Sprintf("%-8s: %s   %s   %s", level, getTimeStamp(), function, message)
	log.Println(fmt.Sprintf(color, logMessage))

	// Logging to a file
	file, err := os.OpenFile("logfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(logMessage)
}

func getTimeStamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
