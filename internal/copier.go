package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// EtherscanMultiFileSource defines the JSON structure for multi-file source code from Etherscan.
type EtherscanMultiFileSource struct {
	Sources map[string]struct {
		Content string `json:"content"`
	} `json:"sources"`
}

// SaveFiles saves a map of source code files to the specified directory.
func SaveFiles(sourceMap map[string]string, outputDir, contractAddress string) error {
	finalOutputDir := filepath.Join(outputDir, contractAddress)
	if err := os.MkdirAll(finalOutputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", finalOutputDir, err)
	}

	fmt.Printf("Found %d files. Saving to directory '%s'...\n", len(sourceMap), finalOutputDir)

	for path, content := range sourceMap {
		fullPath := filepath.Join(finalOutputDir, path)
		fileDir := filepath.Dir(fullPath)

		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create sub-directory '%s': %w", fileDir, err)
		}

		fmt.Printf("- Saving file: %s\n", path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file '%s': %w", path, err)
		}
	}

	fmt.Printf("\n✅ Download complete! Source code saved to directory '%s'.\n", finalOutputDir)
	return nil
}

// ParseEtherscanSource attempts to parse source code from Etherscan, which can be a single string or a JSON object.
func ParseEtherscanSource(sourceCodeJSON string) (map[string]string, error) {
	sourceMap := make(map[string]string)

	// Etherscan might return a JSON string doubly wrapped in braces: {{...}}
	if len(sourceCodeJSON) > 1 && sourceCodeJSON[0] == '{' {
		// Attempt to unmarshal as a multi-file source JSON
		var multiFile EtherscanMultiFileSource

		// Remove the outer braces if they exist
		if sourceCodeJSON[1] == '{' {
			sourceCodeJSON = sourceCodeJSON[1 : len(sourceCodeJSON)-1]
		}

		err := json.Unmarshal([]byte(sourceCodeJSON), &multiFile)
		if err == nil && len(multiFile.Sources) > 0 {
			for path, fileData := range multiFile.Sources {
				sourceMap[path] = fileData.Content
			}
			return sourceMap, nil
		}
	}

	// If it's not a valid multi-file JSON, treat it as a single source file
	sourceMap["contract.sol"] = sourceCodeJSON
	return sourceMap, nil
}
