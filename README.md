# Contract Copier

A command-line tool written in Go to download verified source code for smart contracts from various block explorers.

## Features

- Download verified smart contract source code from Etherscan (for Sepolia Testnet).
- Download verified smart contract source code from zkSync Era Mainnet Explorer.
- Saves source code into organized directories named after the contract address.
- Handles both single-file and multi-file (JSON format) source code structures.

## Installation

### Prerequisites

- Go (version 1.21 or later recommended)
- Git

### From Source

1.  Clone the repository:
    ```sh
    git clone https://github.com/aleister1102/contract-copier.git
    ```
2.  Navigate to the project directory:
    ```sh
    cd contract-copier
    ```
3.  Build the executable:
    ```sh
    go build
    ```
    This will create a `contract-copier` (or `contract-copier.exe` on Windows) executable in the project root.

4.  (Optional) Add the executable to your system's PATH to make it accessible from anywhere. Alternatively, you can run it directly from the project directory (e.g., `./contract-copier` or `.\contract-copier.exe`).

### Using `go install`

If the repository is public and you have Go installed, you can install directly using:

```sh
go install github.com/aleister1102/contract-copier@latest
```

This will install the `contract-copier` binary into your `$GOPATH/bin` or `$HOME/go/bin` directory. Ensure this directory is in your system's PATH.

## Environment Variables

-   `ETHERSCAN_API_KEY`: Required for downloading contracts from Etherscan-powered explorers (like Sepolia). You can get an API key from [etherscan.io](https://etherscan.io/myapikey).

## Usage

The tool uses subcommands to specify the network/explorer to fetch from.

### General Syntax

```sh
contract-copier [command] [contract_address] [flags]
```

### Available Commands

1.  **`sepolia`**: Download contract source code from Sepolia Testnet (via Etherscan).
    ```sh
    contract-copier sepolia <CONTRACT_ADDRESS> -o <OUTPUT_DIRECTORY>
    ```
    -   `<CONTRACT_ADDRESS>`: The address of the smart contract on Sepolia.
    -   `-o, --output <OUTPUT_DIRECTORY>`: (Optional) The directory where the source code will be saved. Defaults to the current directory (`.`). A subdirectory named after the contract address will be created within this output directory.

    **Example:**
    ```sh
    # Make sure ETHERSCAN_API_KEY is set in your environment
    export ETHERSCAN_API_KEY="YOUR_ETHERSCAN_API_KEY"
    contract-copier sepolia 0x123...abc -o ./downloaded_contracts
    ```

2.  **`zksync`**: Download contract source code from zkSync Era Mainnet.
    ```sh
    contract-copier zksync <CONTRACT_ADDRESS> -o <OUTPUT_DIRECTORY>
    ```
    -   `<CONTRACT_ADDRESS>`: The address of the smart contract on zkSync Era Mainnet.
    -   `-o, --output <OUTPUT_DIRECTORY>`: (Optional) The directory where the source code will be saved. Defaults to the current directory (`.`). A subdirectory named after the contract address will be created within this output directory.

    **Example:**
    ```sh
    contract-copier zksync 0x456...def -o ./zk_contracts
    ```

### Global Flags

-   `--help`: Show help information.

## Output

The tool will print the raw API response to the console (for debugging purposes) and then save the extracted source files into a directory structure like:

```
<OUTPUT_DIRECTORY>/
└── <CONTRACT_ADDRESS>/
    ├── Contract1.sol
    ├── Interface.sol
    └── library/
        └── Math.sol
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request or open an Issue.

## License

This project is open-source. Please add a LICENSE file if you wish to specify one (e.g., MIT, Apache 2.0). 
