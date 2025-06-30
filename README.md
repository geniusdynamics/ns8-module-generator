# NS8 Module Generator

## Overview

The NS8 Module Generator is a command-line interface (CLI) tool designed to simplify the creation of NS8 modules from existing Docker Compose files. It automates the process of converting Docker Compose services into NS8-compatible systemd services, handling environment variables, volumes, and other configurations.

## Features

*   **Docker Compose to NS8 Conversion**: Automatically generates NS8 module structure from a Docker Compose file.
*   **Environment Variable Handling**: Supports various Docker Compose environment variable formats.
*   **Volume Management**: Correctly parses and translates Docker Compose volumes into NS8-compatible configurations.
*   **Systemd Service Generation**: Creates systemd service files for each Docker Compose service.
*   **Git Integration**: 
    *   Initializes Git repositories with `main` as the default branch.
    *   Supports both GitHub Tokens and SSH for pushing changes to GitHub.
*   **Template Management**: Downloads and extracts the module template directly into the project root.

## Getting Started

### Prerequisites

*   Go (version 1.23 or higher)
*   Docker Compose file for your application

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/geniusdynamics/ns8-module-generator.git
    cd ns8-module-generator
    ```

2.  **Run the application:**

    ```bash
    go run main.go
    ```

    The CLI will guide you through the process, asking for:
    *   The path to your Docker Compose file.
    *   Your application name.
    *   The output directory for the generated module.
    *   Whether to initialize a Git repository and push to GitHub (if yes, it will ask for your GitHub organization name, username, and preferred authentication method: SSH or Token).

### Building from Source

To build the executable:

```bash
go build -o ns8-module-generator
```

Then, you can run the executable:

```bash
./ns8-module-generator
```

## Usage

Once the tool runs, it will:

1.  Download the latest NS8 module template.
2.  Parse your provided Docker Compose file.
3.  Generate the NS8 module structure in the specified output directory.
4.  Optionally initialize a Git repository and push the generated module to your GitHub account.

## Contributing

We welcome contributions! Please feel free to submit issues or pull requests.

## License

This project is licensed under the [MIT License](LICENSE).