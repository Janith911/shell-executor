package main

import (
	"database/sql"
	"log"
	"os"
	"os/exec"

	"github.com/robfig/cron"
)

func startCron(scripts Scripts, ch chan Scripts, shell string, InfoLogger *log.Logger, ErrLogger *log.Logger, db *sql.DB) {

	c := cron.New()

	for _, v := range scripts.Scripts {
		v := v
		c.AddFunc(v.CronExpression, func() {
			_, err := os.Stat(v.ScriptPath)
			if err != nil {
				ErrLogger.Println("Script ID :", v.Name, "| Error reading sciprt file")
				ErrLogger.Println("Script ID :", v.Name, "| STDERR : ", err)
			} else {
				InfoLogger.Println("Script ID :", v.Name, "| Started Execution")
			}
			cmd := exec.Command(shell, v.ScriptPath)
			// cmd.SysProcAttr = &syscall.SysProcAttr{}
			// cmd.SysProcAttr.Credential = &syscall.Credential{Uid: 501, Gid: 80}
			out, err := cmd.Output()
			if err != nil {
				ErrLogger.Println("Script ID :", v.Name, "| Execution Failed")
				ErrLogger.Println("Script ID :", v.Name, "| STDERR : ", err)
				insertData(db, execution_table_name, v.Name, "now", "FAILED", "SCHEDULED")
			} else {
				InfoLogger.Println("Script ID :", v.Name, "| Execution Successful")
				insertData(db, execution_table_name, v.Name, "now", "SUCCESS", "SCHEDULED")
			}
			if string(out) != "" {
				InfoLogger.Println("Script ID :", v.Name, "STDOUT : ", string(out))
			}

		})
	}

	c.Start()

	ch <- scripts

}
