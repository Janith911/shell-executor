package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func readDbHandler(db *sql.DB, execution_table_name string) func(*gin.Context) {

	return func(c *gin.Context) {
		query := `SELECT * FROM TABLE_NAME ORDER BY id DESC LIMIT 10;`
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

func viewDBOutputClient(url string) {
	resp, e := http.Get(url + "/executions")
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
