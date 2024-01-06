package main

import (
	"log"
	"os"

	"github.com/Flagro/ProjectTextAgent/pkg/watcher"
	"github.com/fsnotify/fsnotify/pkg/database/vecmetaq"
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

	// Create project directory watcher
	w, err := watcher.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}

	err = w.WatchDirectory(projectPath)
	if err != nil {
		log.Fatal("Error watching directory:", err)
	}

	// Watch the project directory for changes and parse the changed file
	for {
		select {
		case filePath := <-w.FileModified:
			log.Println("Modified file:", filePath)
			// parseFile(filePath, tempPath, projectPath, ignorePatterns)
		case filePath := <-w.FileCreated:
			log.Println("Created file:", filePath)
			// parseFile(filePath, tempPath, projectPath, ignorePatterns)
		case filePath := <-w.FileDeleted:
			log.Println("Deleted file:", filePath)
			// handle file deletion...
		}
	}
}
