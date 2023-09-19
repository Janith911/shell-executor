package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/rodaine/table"
)

type DBOutput struct {
	Id          string `json:"id"`
	ScriptName  string `json:"ScriptName"`
	StartTime   string `json:"StartTime"`
	Status      string `json:"Status"`
	ExecutionId string `json:"ExecutionId"`
}

type ManualPayload struct {
	ScriptName string `json:"ScriptName"`
	Shell      string `json:"Shell"`
}

type ManualResponse struct {
	Status string `json:"Status"`
}

// HTTP SERVER FUNCTIONS
func readDbHandler(db *sql.DB, execution_table_name string) func(*gin.Context) {

	return func(c *gin.Context) {
		query := `SELECT * FROM (SELECT * FROM TABLE_NAME ORDER BY id DESC LIMIT 10) ORDER BY id ASC;`
		query = strings.Replace(query, "TABLE_NAME", execution_table_name, 1)
		rows, err := db.Query(query, execution_table_name)
		if err != nil {
			fmt.Println("An ERROR occured : ", err)
		}
		var AllRows []DBOutput

		for rows.Next() {
			var rowOut DBOutput

			rows.Scan(&rowOut.Id, &rowOut.ScriptName, &rowOut.StartTime, &rowOut.Status, &rowOut.ExecutionId)
			AllRows = append(AllRows, rowOut)
		}
		rows.Close()
		c.JSON(http.StatusOK, AllRows)
	}
}

func manualExecutionHandler(configurations Config, InfoLogger *log.Logger, ErrLogger *log.Logger, db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var payload ManualPayload
		var response ManualResponse
		c.BindJSON(&payload)
		out := manualExecution(configurations, payload.ScriptName, payload.Shell, InfoLogger, ErrLogger, db)
		response.Status = out
		c.JSON(http.StatusOK, response)
	}
}

// HTTP CLIENT FUNCTIONS
func viewDBOutputClient(url string) {
	req, err := http.NewRequest("GET", url+"/executions", nil)
	if err != nil {
		fmt.Println("An ERROR Occured : ", err)
	}
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		fmt.Println("An ERROR Occured", e)
	}
	bs, e := io.ReadAll(resp.Body)
	if e != nil {
		fmt.Println("An ERROR Occured : ", e)
	}
	var obj []DBOutput
	json.Unmarshal(bs, &obj)
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "ScriptName", "StartTime", "Status", "ExecutionId")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, e := range obj {
		tbl.AddRow(e.Id, e.ScriptName, e.StartTime, e.Status, e.ExecutionId)
	}
	tbl.Print()
}

func manuallyExecute(url string, scriptName string, shell string) {
	var payload ManualPayload = ManualPayload{
		ScriptName: scriptName,
		Shell:      shell,
	}
	var payloadBuff bytes.Buffer
	err := json.NewEncoder(&payloadBuff).Encode(payload)
	if err != nil {
		fmt.Println("An ERROR Occured : ", err)
	}

	req, e := http.NewRequest("POST", url+"/execute", &payloadBuff)

	if e != nil {
		fmt.Println("An ERROR Occured", e)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("An ERROR Occured", err)
	}
	bs, e := io.ReadAll(resp.Body)
	if e != nil {
		fmt.Println("An ERROR Occured : ", e)
	}
	var obj ManualResponse
	json.Unmarshal(bs, &obj)
	fmt.Println(obj.Status)
}
