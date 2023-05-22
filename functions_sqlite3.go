package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var sql_db_path string = "./test.db"
var execution_table_name string = "executions"

// func main() {
// 	createTable(openDB(sql_db_path), execution_table_name)
// 	fmt.Println(tableExists(openDB(sql_db_path), execution_table_name))
// 	insertData(openDB(sql_db_path), execution_table_name)
// }

func openDB(sql_db_path string) *sql.DB {
	db, err := sql.Open("sqlite3", sql_db_path)

	if err != nil {
		fmt.Println("An ERROR Occured while reading database : ", err)
	}

	return db
}

func createTable(db *sql.DB, execution_table_name string) {
	query := `CREATE TABLE TABLE_NAME (
        executionId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "ScriptName" TEXT,
        "StartTime" TEXT,
        "Status" TEXT,
        "ExecutionType" TEXT);`

	query = strings.Replace(query, "TABLE_NAME", execution_table_name, 1)

	res, err := db.Exec(query)
	if err != nil {
		fmt.Println("An ERROR occured while creating table : ", err)
	}

	n, _ := res.RowsAffected()
	fmt.Println("Rows Affected : ", n)

}
func tableExists(db *sql.DB, execution_table_name string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=$1;`

	rows, err := db.Query(query, execution_table_name)
	if err != nil {
		fmt.Println("An ERROR occured : ", err)
	}
	var table_name string
	for rows.Next() {
		rows.Scan(&table_name)
		if table_name == execution_table_name {
			break
		}
	}

	rows.Close()

	if table_name == execution_table_name {
		return true
	} else {
		return false
	}

}

func insertData(db *sql.DB, execution_table_name string, scriptId string, timeStamp string, status string, executionType string) {
	query := `INSERT INTO TABLE_NAME (ScriptName,StartTime,Status,ExecutionType) VALUES($2,$3,$4,$5)`
	query = strings.Replace(query, "TABLE_NAME", execution_table_name, 1)

	_, err := db.Exec(query, scriptId, timeStamp, status, executionType)

	if err != nil {
		fmt.Println("An ERROR occured while inserting data : ", err)
	}

}
