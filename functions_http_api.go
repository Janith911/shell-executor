package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
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
func readDbHandler(db *sql.DB, execution_table_name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			query := `SELECT * FROM TABLE_NAME ORDER BY id ASC;`
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
			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(AllRows)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, buf.String())
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func manualExecutionHandler(configurations Config, InfoLogger *log.Logger, ErrLogger *log.Logger, db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var payload ManualPayload
			var response ManualResponse
			bs, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("An ERROR Occured : ", err)
			}
			json.Unmarshal(bs, &payload)
			out := manualExecution(configurations, payload.ScriptName, payload.Shell, InfoLogger, ErrLogger, db)
			response.Status = out
			bs, e := json.Marshal(response)
			if e != nil {
				fmt.Println("An ERROR Occured : ", e)
			}
			fmt.Fprintf(w, string(bs))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func listScriptsHandler(configuration Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bs, err := json.Marshal(configuration.Scripts)
			if err != nil {
				fmt.Println("An ERROR Occured : ", err)
			}
			fmt.Fprintf(w, string(bs))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}

func readScriptHandler(configuration Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			scriptId := r.URL.Query().Get("id")
			var scriptPath string
			for _, e := range configuration.Scripts {
				if e.Name == scriptId {
					scriptPath = e.ScriptPath
				}
			}
			if scriptPath == "" {
				fmt.Fprintf(w, "No such script")
			} else {
				_, err := os.Stat(scriptPath)
				if err != nil {
					fmt.Fprintf(w, "No such file in specified path")
				} else {
					bs, e := os.ReadFile(scriptPath)
					if err != nil {
						fmt.Println("An ERROR Occured : ", e)
					}
					fmt.Fprintf(w, string(bs))
				}
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
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

func manuallyExecuteClient(url string, scriptName string, shell string) {
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

func listScriptsClient(url string) {
	req, err := http.NewRequest("GET", url+"/list", nil)
	if err != nil {
		fmt.Println("An ERROR Occured : ", err)
	}
	client := http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		fmt.Println("An ERROR Occured : ", e)
	}
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("An ERROR Occured : ", err)
	}
	var scripts []Script
	err = json.Unmarshal(bs, &scripts)
	if err != nil {
		fmt.Println("An ERROR Occured : ", err)
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Script Name (ID)", "Script Path", "CRON Expression")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, e := range scripts {
		tbl.AddRow(e.Name, e.ScriptPath, e.CronExpression)
	}
	tbl.Print()
}

func readScriptClient(url string, scriptName string) {
	req, e := http.NewRequest("GET", url+"/read?id="+scriptName, nil)
	if e != nil {
		fmt.Println("An ERROR Occured : ", e)
	}
	client := http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		fmt.Println("An ERROR Occured : ", e)
	}
	bs, e := io.ReadAll(resp.Body)
	if e != nil {
		fmt.Println("An ERROR Occured : ", e)
	}
	fmt.Println(string(bs))
}
