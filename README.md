# tiggo

**tiggo** is a command-line tool for devs to supplement their git workflows using LLMs.

## Features

- **Automated Diff Summary Generation:** Automatically generate summaries of code changes to use in PR descriptions.

### Planned Features

- Support for multiple LLM APIs
- More control over the app configs (e.g. prompt templates, output format, etc.)
- Integration with Github CLI for fully automated PR creation

## Project Structure

A simplified overview of the project structure:

```
tiggo/
├── tig.go        # CLI entrypoint
├── internal/
│   ├── analyze/    # Analysis logic (e.g., analyze.go)
│   ├── cli/        # CLI flag parsing and diff range building (e.g., cli.go)
│   ├── git/        # Git integration (e.g., git.go)
│   ├── llm/        # LLM client implementation (e.g., llm.go)
│   └── prompt/     # Prompt templates for LLM requests (e.g., prompt.go)
├── install.sh      # Installation script
├── go.mod          # Go module definition
└── README.md       # This file
```

## Installation

### Prerequisites

- **Go:** Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/doc/install).
- **Git:** The tool relies on git commands, so Git must be installed.

### Building from Source

1. **Clone the repository:**

   ```bash
   git clone https://github.com/steenfuentes/tiggo.git
   cd tiggo
   ```

2. **Download dependencies:**

   ```bash
   go mod tidy
   ```

3. **Install**

   **Install using the provided script:**

   The installation script will build the executable and copy it to a directory in your PATH.

   ```bash
   chmod +x install.sh
   ./install.sh
   ```

   The script checks for Go, builds the executable, and installs it to `/usr/local/bin` (or `$HOME/.local/bin` if needed). After installation, you might need to restart your terminal or run `source ~/.bashrc` or `source ~/.zshrc`.

   **OR**

   **Build the application manually:**

   ```bash
   go build -o tiggo ./tig.go
   ```

   This will create a `tiggo` executable in the current directory.
   If you want to install it globally, you can use the installation
   script, or copy the `tiggo` executable to a directory in your PATH (e.g. `/usr/local/bin`).

## Usage

Once installed, you can use `tiggo` from the command line. Here are some usage examples:

- **Analyze changes between two specific commits:**

  ```bash
  tiggo abc123 def456
  ```

- **Analyze a specified number of commits before a target commit:**

  ```bash
  tiggo def456 -p 5
  ```

- **Analyze the changes from the merge base with main to HEAD (default):**

  ```bash
  tiggo
  ```

For detailed usage information, run:

```bash
tiggo --help
```

## Configuration

The tool first checks for the API key in the system environment. If not found, it will check for a `.env` file in the directory where the command is executed.

Set the `ANTHROPIC_API_KEY` environment variable or create a `.env` file with the following content:

```dotenv
ANTHROPIC_API_KEY=your_api_key_here
```

## Contributing

 Please open an issue or submit a pull request if you have suggestions or improvements.

## License

This project is licensed under the MIT License.

---

This README gives an overview of what tiggo does, details how to build and install it, and provides guidance on usage. The design of the tool—including concurrent file processing, git integration, and LLM API calls—can be explored further in the codebase (see [cli.go](./internal/cli/cli.go), [llm.go](./internal/llm/llm.go), [git.go](./internal/git/git.go), and [analyze.go](./internal/analyze/analyze.go)).