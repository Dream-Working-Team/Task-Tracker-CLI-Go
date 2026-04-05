# Task Tracker CLI (Go) 🚀

![Go Version](https://img.shields.io/badge/Go-1.25-blue)
![Build](https://img.shields.io/badge/CI-GitHub%20Actions-success)
![License](https://img.shields.io/badge/License-MIT-green)

Task Tracker CLI is a multi-user command-line task manager built in Go. It supports user authentication, per-user task storage, task state transitions, and persistent JSON-based data storage.

Current version: `dev` for local builds. Release builds inject the tag version at compile time.

## Author 👤

- Andrés Tortolero
- Instagram: [@andress.tor](https://www.instagram.com/andress.tor?igsh=enRtNXVpNDkxNHNx&utm_source=qr)
- LinkedIn: to be added later

## Overview 📌

The application is organized around a root command powered by [Cobra](https://github.com/spf13/cobra). Each authenticated user gets an isolated task file, so task operations remain scoped to the active session.

## Architecture 🏛️

- **CLI layer (`cmd/`)**: Command routing and UX are implemented with Cobra.
- **Auth layer (`internal/auth/`)**: User registration and login are handled with password hashing via `bcrypt`.
- **Service layer (`internal/service/`)**: Business rules (duplicate validation, status transitions, soft delete visibility) are enforced here.
- **Storage layer (`internal/storage/`)**: Generic JSON file persistence for users and tasks.
- **Isolation model**: Active session is resolved from `.ActiveSession`, and each user writes to `data/task_<user_id>.json`.

This separation keeps CLI concerns decoupled from business logic and persistence, making the project easier to evolve and test.

## Features ✨

- Multi-user authentication with register, login, and logout flows. 👤
- Per-user task storage, isolated by session. 🔐
- Persistent JSON storage for users and tasks. 💾
- CRUD task management: add, update, delete. 🛠️
- Task state changes: `todo`, `in-progress`, and `done`. 🔄
- Filtered task listing by state. 🔎
- Soft delete support so removed tasks do not appear in listings. 🗑️
- Duplicate task description validation to prevent repeated entries. ✅
- Cobra-based CLI structure with built-in help pages per command. 🐍
- Clear command output and error handling for invalid IDs, invalid states, and missing sessions. ⚠️
- Local session tracking to keep the active user between commands. 🧾

## Project Structure 🧱

- `cmd/`: Cobra command definitions and CLI entry points.
- `internal/auth/`: Registration, login, logout, and session handling.
- `internal/config/`: Local data directory resolution.
- `internal/model/`: Domain entities and task status constants.
- `internal/service/`: Business logic for authentication and task management.
- `internal/storage/`: JSON file persistence helpers.
- `data/`: Runtime JSON data files and active session file.

## Command Set ⌨️

### Authentication 🔑

- `task-cli auth register` - Create a new user.
- `task-cli auth login` - Start an authenticated session.
- `task-cli auth logout` - Close the current session.

### Tasks 📋

- `task-cli add "task description"` - Create a new task.
- `task-cli update <id> "new description"` - Update an existing task.
- `task-cli delete <id>` - Soft delete a task.
- `task-cli mark-todo <id>` - Mark a task as to-do.
- `task-cli mark-in-progress <id>` - Mark a task as in progress.
- `task-cli mark-done <id>` - Mark a task as done.
- `task-cli list` - List all active tasks.
- `task-cli list todo|in-progress|done` - Filter tasks by state.

## Data Storage 🗂️

- User records are stored in `data/users.json`.
- Each user has a dedicated task file such as `data/task_<user_id>.json`.
- The active session is stored locally in `data/.ActiveSession`.

This project ignores runtime JSON files through `.gitignore` so local data does not end up in version control.

## Help System ❓

The CLI uses Cobra's help system, so every command exposes its own usage text and examples. You can inspect command-specific help with:

```bash
task-cli --help
task-cli auth --help
task-cli add --help
```

## Installation 🧰

Make sure you have [Go installed](https://go.dev/doc/install).

Clone the repository:

```bash
git clone https://github.com/Dream-Working-Team/Task-Tracker-CLI-Go.git
cd Task-Tracker-CLI-Go
```

Download dependencies:

```bash
go mod download
```

## Build and Run 🏗️

Build the binary:

```bash
go build -o task-cli .
```

Build a release binary with a version string:

```bash
go build -ldflags "-s -w -X task-cli/cmd.Version=v1.0.0" -o task-cli .
```

Run it directly without building:

```bash
go run . --help
```

Check the version from the CLI:

```bash
task-cli --version
```

## CI with GitHub Actions ⚙️

This repository includes a CI workflow at `.github/workflows/go.yml`.

On every push and pull request to `main`, GitHub Actions will:

1. Set up Go.
2. Download dependencies.
3. Build all packages with `go build -v ./...`.
4. Run tests with `go test -v ./...`.

This ensures the project remains buildable and prevents broken changes from being merged.

## Releases and Cross-Compilation 📦

This repository includes an automated release workflow at `.github/workflows/release.yml`.

When you push a tag like `v1.0.0`, GitHub Actions builds binaries for:

- Linux (`amd64`, `arm64`)
- Windows (`amd64`)
- macOS (`arm64`)

Built binaries are attached automatically to the corresponding GitHub Release.

### Manual cross-compilation examples

```bash
# Windows (amd64)
GOOS=windows GOARCH=amd64 go build -o task-cli.exe .

# macOS Apple Silicon (arm64)
GOOS=darwin GOARCH=arm64 go build -o task-cli-mac .

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o task-cli-linux .
```

## Dependabot Security Updates 🔒

Dependabot is configured in `.github/dependabot.yml` for:

- Go modules (`gomod`)
- GitHub Actions dependencies

It checks updates daily and opens pull requests automatically when newer dependency versions are available.

## Usage Example ▶️

```bash
task-cli auth register
task-cli auth login
task-cli add "Study Go"
task-cli mark-in-progress 1
task-cli mark-done 1
task-cli list done
```

## Notes 📝

- Task descriptions are validated to avoid duplicates for the same active user.
- Deleted tasks remain stored but are excluded from normal listings.
- State values are normalized inside the service layer, so the CLI stays consistent across commands.
