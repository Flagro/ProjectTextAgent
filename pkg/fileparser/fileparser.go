package fileparser

import (
	"fmt"
	"log"
	"os/exec"
)

type TextTableScoopOutput struct {
	FilePath string `json:"file_path"`
	Data     []struct {
		DataType string      `json:"data_type"`
		Text     string      `json:"text"`
		Metadata interface{} `json:"metadata"`
	} `json:"data"`
}

func ParseFile(filePath, tempPath, projectPath, ignorePatterns string) string {
	// Construct the command
	cmd := exec.Command("texttablescoop", filePath, "--temp", tempPath, "--project", projectPath, "--ignore", ignorePatterns)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	fmt.Println(string(output))
	return string(output)
}
