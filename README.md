# copr2url

`copr2url` is a Go-based command-line tool to fetch direct RPM links from Fedora's Copr repositories. It retrieves the latest successful build for specified packages and generates download links for `x86_64` or `noarch` architectures.

## Features

- Retrieves latest successful builds via Copr API.
- Extracts RPM links using Copr HTML directory parsing.
- Supports plain and JSON output formats.

## Usage

### CLI Syntax

```bash
copr2url [repos.ini] [fedora-target] [--json]
```

- `repos.ini`: Path to the INI file containing repository details (default: `repos.ini`).
- `fedora-target`: Fedora build target (e.g., `fedora-41-x86_64`, `fedora-rawhide-x86_64`; default: `fedora-rawhide-x86_64`).
- `--json`: Outputs results in JSON format.

### Examples

1. Generate plain RPM links:
   ```bash
   copr2url
   ```
2. Use custom INI file and target:
   ```bash
   copr2url myrepos.ini fedora-41-x86_64
   ```
3. Output RPM links in JSON format:
   ```bash
   copr2url repos.ini fedora-41-x86_64 --json
   ```

## Input Format

### INI File Example

```ini
[cosmic-app-library]
owner = ryanabx
project = cosmic-epoch
package = cosmic-app-library

[cosmic-applets]
owner = ryanabx
project = cosmic-epoch
package = cosmic-applets
```

## Build & Run

1. Install Go.
2. Clone the repository.
3. Build the binary:
   ```bash
   go build -o copr2url main.go
   ```
4. Run:
   ```bash
   ./copr2url [repos.ini] [fedora-target] [--json]
   ```
