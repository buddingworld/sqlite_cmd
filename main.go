package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: program <db_file> <select|update> <sql_query>")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	queryType := os.Args[2]
	sqlStr := os.Args[3]

	if queryType != "select" && queryType != "update" {
		fmt.Println("Query type must be 'select' or 'update'")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	switch queryType {
	case "select":
		handleSelect(db, sqlStr)
	case "update":
		handleUpdate(db, sqlStr)
	}
}

func handleSelect(db *sql.DB, sqlStr string) {
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("Failed to get columns: %v\n", err)
		os.Exit(1)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			fmt.Printf("Row scan failed: %v\n", err)
			os.Exit(1)
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			entry[col] = values[i]
		}
		results = append(results, entry)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Row iteration error: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("JSON marshal failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonData))
}

func handleUpdate(db *sql.DB, sqlStr string) {
	result, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("Update failed: %v\n", err)
		os.Exit(1)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Failed to get rows affected: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Rows affected: %d\n", rowsAffected)
}