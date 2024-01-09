package main

import (
	"log"
	"os"

	"github.com/Flagro/ProjectTextAgent/pkg/database/postgres"
	"github.com/Flagro/ProjectTextAgent/pkg/database/vecmetaq"
	"github.com/Flagro/ProjectTextAgent/pkg/fileparser"
	"github.com/Flagro/ProjectTextAgent/pkg/watcher"
)

func updateDataBases(postgresClient *postgres.Client, vecmetaqClient *vecmetaq.Client, parser_output *fileparser.TextTableScoopOutput) {

}

func removeFromDatabases(postgresClient *postgres.Client, vecmetaqClient *vecmetaq.Client, filePath string) {

}

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
	vecmetaqClient, err := vecmetaq.NewClient(vecmetaqURL, vecmetaqUsername, vecmetaqPassword)
	if err != nil {
		log.Fatal("Error creating VecMetaQ client:", err)
	}
	defer vecmetaqClient.Close()

	postgresClient, err := postgres.NewClient(postgresURL, postgresUser, postgresPassword)
	if err != nil {
		log.Fatal("Error creating PostgreSQL client:", err)
	}
	defer postgresClient.Close()

	// if both dbs are empty, parse the whole project directory
	if postgresClient.IsEmpty() && vecmetaqClient.IsEmpty() {
		wholeProjectOutput := fileparser.ParseFile(projectPath, tempPath, projectPath, ignorePatterns)
		updateDataBases(postgresClient, vecmetaqClient, &wholeProjectOutput)
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
			fileOutput := fileparser.ParseFile(filePath, tempPath, projectPath, ignorePatterns)
			updateDataBases(postgresClient, vecmetaqClient, &fileOutput)
		case filePath := <-w.FileCreated:
			log.Println("Created file:", filePath)
			fileOutput := fileparser.ParseFile(filePath, tempPath, projectPath, ignorePatterns)
			updateDataBases(postgresClient, vecmetaqClient, &fileOutput)
		case filePath := <-w.FileDeleted:
			log.Println("Deleted file:", filePath)
			removeFromDatabases(postgresClient, vecmetaqClient, filePath)
		}
	}
}
