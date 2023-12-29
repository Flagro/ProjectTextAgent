package main

import (
	"log"
)

func main() {
	path := "/path/to/your/directory"
	tempPath := "path/to/temp"
	projectPath := "path/to/project"
	ignorePatterns := "pattern1,pattern2"

	// Initialize PostgreSQL connection
	pgDB, err := postgres.New("your_postgres_connection_string")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDB.Close()

	// Initialize VectorDB connection
	vectorDB, err := vectordb.New(yourVectorDBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to VectorDB: %v", err)
	}

	fileChanged := make(chan string)
	go watchDirectory(path, fileChanged)

	for {
		select {
		case filePath := <-fileChanged:
			log.Println("Need to handle file:", filePath)
			parseFile(filePath, tempPath, projectPath, ignorePatterns)
			// updateDatabase(parsedData)
		}
	}
}
