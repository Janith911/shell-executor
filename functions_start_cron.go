package main

import (
	"database/sql"
	"log"
	"os"
	"os/exec"
	"time"

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
				ErrLogger.Println("Script ID : ", v.Name, "| Execution Failed")
				ErrLogger.Println("Script ID : ", v.Name, "| STDERR : \n", err)
				insertData(db, execution_table_name, v.Name, time.Now().Format(time.RFC3339), "FAILED", "Scheduled_"+time.Now().Format("2006_01_02_15:04:05"))
			} else {
				InfoLogger.Println("Script ID : ", v.Name, "| Execution Successful")
				insertData(db, execution_table_name, v.Name, time.Now().Format(time.RFC3339), "SUCCESS", "Scheduled_"+time.Now().Format("2006_01_02_15:04:05"))
			}
			if string(out) != "" {
				InfoLogger.Println("Script ID : ", v.Name, "STDOUT : \n", string(out))
			}

		})
	}

	c.Start()

	ch <- scripts

}
