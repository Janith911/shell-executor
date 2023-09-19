package main

import (
	"log"
	"net/http"
	"os"
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
		// START API
		// fmt.Println("Starting HTTP Endpoint")
		log.Println("Starting HTTP Endpoint")
		http.HandleFunc("/executions", readDbHandler(db, execution_table_name))
		http.HandleFunc("/execute", manualExecutionHandler(jsonConf, &InfoLogger, &ErrLogger, db))
		go http.ListenAndServe(jsonConf.BindIP+":"+jsonConf.BindPort, nil)
		// fmt.Println("Started HTTP Endpoint successfully")
		log.Println("Started HTTP Endpoint successfully")
		log.Println("Listening on : http://" + jsonConf.BindIP + ":" + jsonConf.BindPort)

		// START SCHEDULER
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

	} else if os.Args[1] == "execute" {
		if len(os.Args) == 4 {
			url := "http://127.0.0.1:8080"
			scriptName := os.Args[2]
			shell := os.Args[3]
			manuallyExecute(url, scriptName, shell)
		}
	}
}
