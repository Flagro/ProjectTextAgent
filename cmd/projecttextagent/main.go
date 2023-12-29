package main

import (
	"log"
)

func main() {
	path := "/path/to/your/directory"
	tempPath := "path/to/temp"
	projectPath := "path/to/project"
	ignorePatterns := "pattern1,pattern2"
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
