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
	for _, entry := range *parser_output {
		filePath := entry.FilePath
		postgresClient.RemoveData(filePath)
		vecmetaqClient.RemoveData(filePath)
		for _, data := range entry.Data {
			dataType := data.DataType
			text := data.Text
			metadata := string(data.Metadata)
			if dataType == "table" {
				err := postgresClient.AddData(filePath, text, metadata)
				if err != nil {
					log.Fatal("Error adding data to PostgreSQL:", err)
				}
			} else if dataType == "text" {
				err := vecmetaqClient.AddData(filePath, text, metadata)
				if err != nil {
					log.Fatal("Error adding data to VecMetaQ:", err)
				}
			}
		}
	}
}

func removeFromDatabases(postgresClient *postgres.Client, vecmetaqClient *vecmetaq.Client, filePath string) {
	err := postgresClient.RemoveData(filePath)
	if err != nil {
		log.Fatal("Error removing data from PostgreSQL:", err)
	}
	err = vecmetaqClient.RemoveData(filePath)
	if err != nil {
		log.Fatal("Error removing data from VecMetaQ:", err)
	}
}

func main() {
	// ProjectTextAgent configuration
	projectPath := os.Getenv("PROJECT_PATH")
	tempPath := os.Getenv("TEMP_FOLDER_PATH")
	ignorePatterns := os.Getenv("IGNORE_PATTERS")

	// VecMetaQ configuration
	vecmetaqHost := os.Getenv("VECMETAQ_HOST")
	vecmetaqPort := os.Getenv("VECMETAQ_PORT")
	vecmetaqUsername := os.Getenv("VECMETAQ_USER")
	vecmetaqPassword := os.Getenv("VECMETAQ_PASSWORD")

	// PostgreSQL configuration
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_NAME")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")

	// Initialize the databases connections
	vecmetaqClient, err := vecmetaq.NewClient(vecmetaqHost, vecmetaqPort, vecmetaqUsername, vecmetaqPassword)
	if err != nil {
		log.Fatal("Error creating VecMetaQ client:", err)
	}
	defer vecmetaqClient.Close()

	postgresClient, err := postgres.NewClient(postgresHost, postgresPort, postgresDB, postgresUser, postgresPassword)
	if err != nil {
		log.Fatal("Error creating PostgreSQL client:", err)
	}
	defer postgresClient.Close()

	// if both dbs are empty, parse the whole project directory
	postgresIsEmpty, err := postgresClient.IsEmpty()
	if err != nil {
		log.Fatal("Error checking if PostgreSQL database is empty:", err)
	}

	vecmetaqIsEmpty, err := vecmetaqClient.IsEmpty()
	if err != nil {
		log.Fatal("Error checking if VecMetaQ database is empty:", err)
	}

	if postgresIsEmpty && vecmetaqIsEmpty {
		wholeProjectOutput, err := fileparser.ParseFile(projectPath, tempPath, projectPath, ignorePatterns)
		if err != nil {
			log.Println("Error parsing project directory:", err)
		} else {
			updateDataBases(postgresClient, vecmetaqClient, &wholeProjectOutput)
		}
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
			fileOutput, err := fileparser.ParseFile(filePath, tempPath, projectPath, ignorePatterns)
			if err != nil {
				log.Println("Error parsing file:", err)
				continue
			}
			updateDataBases(postgresClient, vecmetaqClient, &fileOutput)
		case filePath := <-w.FileCreated:
			log.Println("Created file:", filePath)
			fileOutput, err := fileparser.ParseFile(filePath, tempPath, projectPath, ignorePatterns)
			if err != nil {
				log.Println("Error parsing file:", err)
				continue
			}
			updateDataBases(postgresClient, vecmetaqClient, &fileOutput)
		case filePath := <-w.FileDeleted:
			log.Println("Deleted file:", filePath)
			removeFromDatabases(postgresClient, vecmetaqClient, filePath)
		}
	}
}
