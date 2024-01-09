package fileparser

import (
	"encoding/json"
	"log"
	"os/exec"
)

type TextTableScoopOutput []struct {
	FilePath string `json:"file_path"`
	Data     []struct {
		DataType string      `json:"data_type"`
		Text     string      `json:"text"`
		Metadata interface{} `json:"metadata"`
	} `json:"data"`
}

func ParseFile(filePath, tempPath, projectPath, ignorePatterns string) TextTableScoopOutput {
	// Construct the command
	cmd := exec.Command("texttablescoop", filePath, "--temp", tempPath, "--project", projectPath, "--ignore", ignorePatterns)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("texttablescoop finished with error: %v", err)
	}

	// Unmarshal the output into a TextTableScoopOutput
	var result TextTableScoopOutput
	err = json.Unmarshal(output, &result)
	if err != nil {
		log.Fatalf("Failed to unmarshal output: %v", err)
	}

	return result
}
