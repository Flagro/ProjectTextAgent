package main

import (
	"database/sql"
	"log"

	"github.com/fsnotify/fsnotify"
)

// TODO: add recursive watcher
func watchDirectory(path string, fileChanged chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Process events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file:", event.Name)
					fileChanged <- event.Name // Send the changed file path to the channel
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-fileChanged // Block until a file is changed
}

type ParsedData struct {
	FilePath string `json:"file_path"`
	Data     []struct {
		DataType string                 `json:"data_type"`
		CSVText  string                 `json:"csv_text"`
		Text     string                 `json:"text"`
		Metadata map[string]interface{} `json:"metadata"`
	} `json:"data"`
}

func processData(data []ParsedData, db *sql.DB) {
	for _, entry := range data {
		if hasTableData(entry) {
			removeExistingTables(entry.FilePath, db)
			for _, dataEntry := range entry.Data {
				if dataEntry.DataType == "table" {
					createTableFromCSV(dataEntry.CSVText, entry.FilePath, db)
				}
			}
		}
	}
}

func removeExistingTables(filePath string, db *sql.DB) {
	// Implementation depends on how you're tracking the association between file paths and tables
	// Example: DELETE FROM metadata_table WHERE file_path = filePath
}

func createTableFromCSV(csvText, filePath string, db *sql.DB) {
	// Parse the CSV data
	// Create a CREATE TABLE SQL statement based on the CSV structure
	// Execute the SQL statement to create the table
	// Optionally, update a metadata table to track the association between the table and the file path
}
