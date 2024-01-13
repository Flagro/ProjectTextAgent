package fileparser

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type TextTableScoopOutput []struct {
	FilePath string `json:"file_path"`
	Data     []struct {
		DataType string          `json:"data_type"`
		Text     string          `json:"text"`
		Metadata json.RawMessage `json:"metadata"`
	} `json:"data"`
}

func ParseFile(filePath, tempPath, projectPath, ignorePatterns string) (TextTableScoopOutput, error) {
	// Construct the command
	cmd := exec.Command("texttablescoop", filePath, "--temp", tempPath, "--project", projectPath, "--ignore", ignorePatterns)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return TextTableScoopOutput{}, fmt.Errorf("texttablescoop finished with error: %v", err)
	}

	// Unmarshal the output into a TextTableScoopOutput
	var result TextTableScoopOutput
	err = json.Unmarshal(output, &result)
	if err != nil {
		return TextTableScoopOutput{}, fmt.Errorf("failed to unmarshal output: %v", err)
	}

	return result, nil
}
