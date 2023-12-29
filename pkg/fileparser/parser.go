package main

import (
	"fmt"
	"log"
	"os/exec"
)

func parseFile(filePath, tempPath, projectPath, ignorePatterns string) {
	// Construct the command
	cmd := exec.Command("texttablescoop", filePath, "--temp", tempPath, "--project", projectPath, "--ignore", ignorePatterns)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	fmt.Println(string(output))
}
