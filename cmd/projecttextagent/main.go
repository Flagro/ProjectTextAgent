package main

import (
	"log"
	"os"

	"github.com/Flagro/ProjectTextAgent/pkg/database/postgres"
	"github.com/Flagro/ProjectTextAgent/pkg/database/vecmetaq"
)

func main() {
	projectPath := "path/to/project"
	tempPath := "path/to/temp"
	ignorePatterns := "pattern1,pattern2"

	// Initialize PostgreSQL connection
	pgDB, err := postgres.New("your_postgres_connection_string")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDB.Close()

	// Initialize VectorDB connection
	vecmetaqBaseURL := os.Getenv("VECMETAQ_BASE_URL")
	vecmetaqUsername := os.Getenv("VECMETAQ_USERNAME")
	vecmetaqPassword := os.Getenv("VECMETAQ_PASSWORD")

	vecmetaqClient := vecmetaq.NewClient(vecmetaqBaseURL, vecmetaqUsername, vecmetaqPassword)

	// if both dbs are empty, parse the whole project directory
	if pgDB.IsEmpty() && vecmetaqClient.IsEmpty() {
		parseFile(projectPath, tempPath, projectPath, ignorePatterns)
	}

	// Watch the project directory for changes and parse the changed file
	fileModified := make(chan string)
	fileCreated := make(chan string)
	fileDeleted := make(chan string)
	go watchDirectory(path, fileModified, fileCreated, fileDeleted)

	for {
		select {
		case filePath := <-fileChanged:
			log.Println("Need to handle file:", filePath)
			parseFile(filePath, tempPath, projectPath, ignorePatterns)
			// updateDatabase(parsedData)
		case filePath := <-fileCreated:
			log.Println("Need to handle file:", filePath)
			parseFile(filePath, tempPath, projectPath, ignorePatterns)
			// updateDatabase(parsedData)
		case filePath := <-fileDeleted:
			log.Println("Need to handle file:", filePath)
			// updateDatabase(parsedData)
		}
	}
}
