package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var CronEx = make(map[string][]string)

func main() {
	if os.Args[1] == "start" {
		jsonConf := readConfig(os.Getenv("CONFIG_FILE_PATH"))
		var LOG_FILE string = jsonConf.LogFilePath

		var InfoLogger log.Logger = *log.New(openLogFile(LOG_FILE), "[INFO] ", log.Flags())
		var ErrLogger log.Logger = *log.New(openLogFile(LOG_FILE), "[ERROR] ", log.Flags())

		ch := make(chan Config)

		db := openDB(jsonConf.DbFilePath)
		if !tableExists(db, execution_table_name) {
			createTable(db, execution_table_name)
		}

		router := gin.Default()
		router.GET("/executions", readDbHandler(db, execution_table_name))
		go router.Run(jsonConf.BindIP + ":" + jsonConf.BindPort)

		startCron(jsonConf, ch, "/bin/bash", &InfoLogger, &ErrLogger, db)

		for s := range ch {
			startCron(s, ch, "/bin/bash", &InfoLogger, &ErrLogger, db)
		}
	} else if os.Args[1] == "executions" {
		if len(os.Args) == 2 {
			viewDBOutputClient("http://127.0.0.1:8080")
		} else if len(os.Args) == 3 {
			viewDBOutputClient(os.Args[2])
		}

	}
}
