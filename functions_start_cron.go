package main

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron"
)

func startCron(scripts Config, ch chan Config, shell string, InfoLogger *log.Logger, ErrLogger *log.Logger, db *sql.DB) {

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
			var cmd_out bytes.Buffer
			var cmd_err bytes.Buffer
			cmd.Stdout = &cmd_out
			cmd.Stderr = &cmd_err
			err = cmd.Run()
			if err != nil {
				ErrLogger.Println("Script ID : ", v.Name, "| Execution Failed")
				ErrLogger.Println("Script ID : ", v.Name, "| STDERR : \n", cmd_err.String())
				insertData(db, execution_table_name, v.Name, time.Now().Format(time.RFC3339), "FAILED", "Scheduled_"+time.Now().Format("2006_01_02_15:04:05"))
			} else {
				InfoLogger.Println("Script ID : ", v.Name, "| Execution Successful")
				InfoLogger.Println("Script ID : ", v.Name, "STDOUT : \n", string(cmd_out.String()))
				insertData(db, execution_table_name, v.Name, time.Now().Format(time.RFC3339), "SUCCESS", "Scheduled_"+time.Now().Format("2006_01_02_15:04:05"))
			}

		})
	}

	c.Start()

	ch <- scripts

}

func manualExecution(configurations Config, scriptName string, shell string, InfoLogger *log.Logger, ErrLogger *log.Logger, db *sql.DB) string {
	var scriptDetails Script
	for _, v := range configurations.Scripts {
		if v.Name == scriptName {
			scriptDetails = v
			break
		}
	}
	if scriptDetails.Name == "" {
		return "No Such Script"
	} else {
		_, err := os.Stat(scriptDetails.ScriptPath)
		if err != nil {
			ErrLogger.Println("Script ID :", scriptDetails.Name, "| Error reading sciprt file")
			ErrLogger.Println("Script ID :", scriptDetails.Name, "| STDERR : ", err)
		} else {
			InfoLogger.Println("Script ID :", scriptDetails.Name, "| Started Execution")
		}
		cmd := exec.Command(shell, scriptDetails.ScriptPath)
		var cmd_out bytes.Buffer
		var cmd_err bytes.Buffer
		cmd.Stdout = &cmd_out
		cmd.Stderr = &cmd_err
		err = cmd.Run()
		if err != nil {
			ErrLogger.Println("Script ID : ", scriptDetails.Name, "| Execution Failed")
			ErrLogger.Println("Script ID : ", scriptDetails.Name, "| STDERR : \n", cmd_err.String())
			insertData(db, execution_table_name, scriptDetails.Name, time.Now().Format(time.RFC3339), "FAILED", "Manual_"+time.Now().Format("2006_01_02_15:04:05"))
			return "Execution Failed"
		} else {
			InfoLogger.Println("Script ID : ", scriptDetails.Name, "| Execution Successful")
			InfoLogger.Println("Script ID : ", scriptDetails.Name, "STDOUT : \n", string(cmd_out.String()))
			insertData(db, execution_table_name, scriptDetails.Name, time.Now().Format(time.RFC3339), "SUCCESS", "Manual_"+time.Now().Format("2006_01_02_15:04:05"))
			return "Execution Success"
		}
	}
}
