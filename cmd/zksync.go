package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aleister1102/contract-cloner/internal"

	"github.com/spf13/cobra"
)

var outputDirZksync string

var zksyncCmd = &cobra.Command{
	Use:   "zksync [contract_address]",
	Short: "Download contract source code from zkSync Era Mainnet",
	Args:  cobra.ExactArgs(1), // Requires exactly one argument
	Run: func(cmd *cobra.Command, args []string) {
		contractAddress := args[0]
		fmt.Printf("Fetching source code for zkSync contract: %s\n", contractAddress)

		apiURL := fmt.Sprintf("https://block-explorer-api.mainnet.zksync.io/api?module=contract&action=getsourcecode&address=%s", contractAddress)

		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Printf("Error calling API: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read the body for printing
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		// fmt.Printf("Raw API Response:\n%s\n", string(bodyBytes))

		// Re-assign resp.Body with a new reader for json.NewDecoder
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Define the structure for the outer API response (like sepolia.go)
		var apiResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Result  []struct {
				SourceCode string `json:"SourceCode"` // This is the JSON string
				// ContractName string `json:"ContractName"` // Optional: could be used later
			} `json:"result"`
		}

		// Decode the main HTTP response body
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error parsing API JSON response: %v\n", err)
			return
		}

		if apiResp.Status != "1" || len(apiResp.Result) == 0 || apiResp.Result[0].SourceCode == "" {
			errorMsg := "Failed to retrieve source code."
			if apiResp.Message != "" {
				errorMsg = fmt.Sprintf("Failed to retrieve source code: %s (Status: %s)", apiResp.Message, apiResp.Status)
			} else if apiResp.Status != "1" {
				errorMsg = fmt.Sprintf("API request failed with status: %s", apiResp.Status)
			} else if len(apiResp.Result) == 0 {
				errorMsg = "API request successful, but no results returned."
			} else { // SourceCode string is empty
				errorMsg = "API request successful, result found, but source code string is empty."
			}
			fmt.Println(errorMsg)
			return
		}

		sourceCodeString := apiResp.Result[0].SourceCode

		// Use internal.ParseEtherscanSource to parse the source code string
		filesToSave, err := internal.ParseEtherscanSource(sourceCodeString)
		if err != nil {
			fmt.Printf("Error processing source code string: %v\n", err)
			// For debugging, you might want to print the sourceCodeString here:
			// fmt.Printf("Problematic SourceCode string: %s\n", sourceCodeString)
			return
		}

		if len(filesToSave) == 0 {
			fmt.Println("No source files extracted after parsing. The SourceCode string might be in an unexpected format or empty.")
			// For debugging, print the sourceCodeString that resulted in empty filesToSave:
			// fmt.Printf("SourceCode string that led to no files: %s\n", sourceCodeString)
			return
		}

		if err := internal.SaveFiles(filesToSave, outputDirZksync, contractAddress); err != nil {
			fmt.Printf("Error saving files: %v\n", err)
		} else {
			fmt.Println("Source code downloaded and saved successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(zksyncCmd)
	zksyncCmd.Flags().StringVarP(&outputDirZksync, "output", "o", ".", "Output directory to save source code")
}
