package main

import (
	"encoding/csv"
	"strings"

	"gorm.io/gorm"
)

type ParsedData struct {
	FilePath string      `json:"file_path"`
	Data     []DataEntry `json:"data"`
}

type DataEntry struct {
	DataType string                 `json:"data_type"`
	CSVText  string                 `json:"csv_text"`
	Text     string                 `json:"text"`
	Metadata map[string]interface{} `json:"metadata"`
}

type TableMetadata struct {
	gorm.Model
	FilePath  string
	TableName string
}

func processData(data []ParsedData, db *gorm.DB) error {
	for _, entry := range data {
		if hasTableData(entry) {
			err := removeExistingTables(entry.FilePath, db)
			if err != nil {
				return err
			}

			for _, dataEntry := range entry.Data {
				if dataEntry.DataType == "table" {
					err := createTableFromCSV(dataEntry.CSVText, entry.FilePath, db)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func hasTableData(entry ParsedData) bool {
	for _, data := range entry.Data {
		if data.DataType == "table" {
			return true
		}
	}
	return false
}

func removeExistingTables(filePath string, db *gorm.DB) error {
	var tables []TableMetadata
	result := db.Where("file_path = ?", filePath).Find(&tables)
	if result.Error != nil {
		return result.Error
	}

	for _, table := range tables {
		db.Exec("DROP TABLE IF EXISTS " + table.TableName)
		db.Delete(&table)
	}
	return nil
}

func createTableFromCSV(csvText, filePath string, db *gorm.DB) error {
	// Parse the CSV data
	reader := csv.NewReader(strings.NewReader(csvText))
	_, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Generate a CREATE TABLE statement (this is an example and needs to be adapted)
	tableName := generateTableName(filePath)
	createStatement := "CREATE TABLE " + tableName + " (id SERIAL PRIMARY KEY, ...);"
	db.Exec(createStatement)

	// Insert metadata entry
	db.Create(&TableMetadata{FilePath: filePath, TableName: tableName})

	return nil
}

func generateTableName(filePath string) string {
	// Implement a method to generate a unique table name based on the file path
	return filePath
}
