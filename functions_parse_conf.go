package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Script struct {
	Name           string
	ScriptPath     string
	CronExpression string
}

type Scripts struct {
	Scripts     []Script
	LogFilePath string
	DbFilePath  string
}

func readConfig(confPath string) string {
	confFile, err := os.Open(confPath)
	if err != nil {
		fmt.Println("Failed Reading Configuration File")
	}

	var b bytes.Buffer
	mw := io.Writer(&b)
	io.Copy(mw, confFile)
	return b.String()

}

func unmarshalConfig(conf string) Scripts {
	bs := []byte(conf)
	scripts := Scripts{}

	err := json.Unmarshal(bs, &scripts)

	if err != nil {
		fmt.Println("Failed Parsing Configuration File : ", err)
		os.Exit(1)
	}

	return scripts
}
