package main

import (
	"log"
)

var CronEx = make(map[string][]string)

func main() {

	jsonConf := unmarshalConfig(readConfig("./conf.json"))

	var LOG_FILE string = jsonConf.LogFilePath

	var InfoLogger log.Logger = *log.New(openLogFile(LOG_FILE), "[INFO] ", log.Flags())
	var ErrLogger log.Logger = *log.New(openLogFile(LOG_FILE), "[ERROR] ", log.Flags())

	ch := make(chan Scripts)

	db := openDB(jsonConf.DbFilePath)
	if !tableExists(db, execution_table_name) {
		createTable(db, execution_table_name)
	}

	startCron(jsonConf, ch, "/bin/bash", &InfoLogger, &ErrLogger, db)

	for s := range ch {
		startCron(s, ch, "/bin/bash", &InfoLogger, &ErrLogger, db)
	}
}
