# Task Tracker CLI (Go) 🚀

Task Tracker CLI is a multi-user command-line task manager built in Go. It supports user authentication, per-user task storage, task state transitions, and persistent JSON-based data storage.

## Overview 📌

The application is organized around a root command powered by [Cobra](https://github.com/spf13/cobra). Each authenticated user gets an isolated task file, so task operations remain scoped to the active session.

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

Run it directly without building:

```bash
go run . --help
```

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
