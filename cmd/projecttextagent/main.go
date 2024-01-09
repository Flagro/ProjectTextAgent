package main

import (
	"log"
	"os"

	"github.com/Flagro/ProjectTextAgent/pkg/database/postgres"
	"github.com/Flagro/ProjectTextAgent/pkg/database/vecmetaq"
	"github.com/Flagro/ProjectTextAgent/pkg/watcher"
)

func main() {
	// ProjectTextAgent configuration
	projectPath := os.Getenv("PROJECT_PATH")
	tempPath := os.Getenv("TEMP_FOLDER_PATH")
	ignorePatterns := os.Getenv("IGNORE_PATTERS")

	// VecMetaQ configuration
	vecmetaqBaseURL := os.Getenv("HOST") // TODO: change to VECMETAQ_HOST
	vecmetaqPort := os.Getenv("PORT")    // TODO: change to VECMETAQ_PORT
	vecmetaqURL := vecmetaqBaseURL + ":" + vecmetaqPort
	vecmetaqUsername := os.Getenv("VECMETAQ_USER")
	vecmetaqPassword := os.Getenv("VECMETAQ_PASSWORD")

	// PostgreSQL configuration
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_NAME")
	postgresURL := postgresHost + ":" + postgresPort + "/" + postgresDB
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	// Initialize the databases connections
	vecmetaqClient := vecmetaq.NewClient(vecmetaqURL, vecmetaqUsername, vecmetaqPassword)
	defer vecmetaqClient.Close()
	postgresClient := postgres.NewClient(postgresURL, postgresUser, postgresPassword)
	defer postgresClient.Close()

	// if both dbs are empty, parse the whole project directory
	if postgresClient.IsEmpty() && vecmetaqClient.IsEmpty() {
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
