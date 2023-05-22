package main

import(
	"os"
	"fmt"
)

func openLogFile(logFilePath string) *os.File {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Print("Error opening logfile : ", err)
		os.Exit(1)
	}
	return logFile
}