package main

import (
	"log"
)

func main() {
	path := "/path/to/your/directory"
	fileChanged := make(chan string)
	go watchDirectory(path, fileChanged)

	for {
		select {
		case filePath := <-fileChanged:
			log.Println("Need to handle file:", filePath)
			// Here you can call the parser and database update logic
			// parseFile(filePath)
			// updateDatabase(parsedData)
		}
	}
}
