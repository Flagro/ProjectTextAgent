package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify/pkg/database/vecmetaq"
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
	vecmetaqBaseURL := os.Getenv("VECMETAQ_BASE_URL")
	vecmetaqUsername := os.Getenv("VECMETAQ_USERNAME")
	vecmetaqPassword := os.Getenv("VECMETAQ_PASSWORD")

	vecmetaqClient := vecmetaq.NewClient(vecmetaqBaseURL, vecmetaqUsername, vecmetaqPassword)

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
