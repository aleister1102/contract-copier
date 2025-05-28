package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"photocopier/internal"

	"github.com/spf13/cobra"
)

var outputDirSepolia string

var sepoliaCmd = &cobra.Command{
	Use:   "sepolia [contract_address]",
	Short: "Download contract source code from Sepolia Testnet (via Etherscan)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		contractAddress := args[0]

		// Read API key from environment variable
		apiKey := os.Getenv("ETHERSCAN_API_KEY")
		if apiKey == "" {
			fmt.Println("Error: ETHERSCAN_API_KEY environment variable is not set.")
			fmt.Println("Please get an API key from etherscan.io and set the environment variable.")
			return
		}

		fmt.Printf("Fetching source code for Sepolia contract: %s\n", contractAddress)
		apiURL := fmt.Sprintf("https://api-sepolia.etherscan.io/api?module=contract&action=getsourcecode&address=%s&apikey=%s", contractAddress, apiKey)

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
		fmt.Printf("Raw API Response:\n%s\n", string(bodyBytes))

		// Re-assign resp.Body with a new reader for json.NewDecoder
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var result struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Result  []struct {
				SourceCode string `json:"SourceCode"`
			} `json:"result"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("Error parsing JSON response: %v\n", err)
			return
		}

		if result.Status != "1" || len(result.Result) == 0 || result.Result[0].SourceCode == "" {
			fmt.Printf("Failed to retrieve source code: %s\n", result.Message)
			return
		}

		sourceCode := result.Result[0].SourceCode
		sourceMap, err := internal.ParseEtherscanSource(sourceCode)
		if err != nil {
			fmt.Printf("Error processing source code: %v\n", err)
			return
		}

		if err := internal.SaveFiles(sourceMap, outputDirSepolia, contractAddress); err != nil {
			fmt.Printf("Error saving files: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sepoliaCmd)
	sepoliaCmd.Flags().StringVarP(&outputDirSepolia, "output", "o", ".", "Output directory to save source code")
}
