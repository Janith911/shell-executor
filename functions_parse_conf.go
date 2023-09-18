package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Script struct {
	Name           string `json:"Name"`
	ScriptPath     string `json:"ScriptPath"`
	CronExpression string `json:"CronExpression"`
}

type Config struct {
	Scripts     []Script `json:"Scripts"`
	LogFilePath string   `json:"LogFilePath"`
	DbFilePath  string   `json:"DbFilePath"`
	BindIP      string   `json:"BindIP"`
	BindPort    string   `json:"BindPort"`
}

func readConfig(confPath string) Config {
	confFile, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Println("Failed Reading Configuration File")
	}

	obj := Config{}
	e := json.Unmarshal(confFile, &obj)

	if e != nil {
		fmt.Println("Failed Unmarshalling Configuration File")
	}

	return obj

}
