# Pygo

Pygo is a lightweight implementation of Git written in Go. It provides basic Git functionality, including repository initialization, object hashing, and inspecting object contents.

## Features

- **Repository Management**: Initialize new Git repositories.
- **Git Objects**: Support for core Git objects:
  - `blob`: Stores file content.
  - `tree`: Stores directory structures.
  - `commit`: Stores commit metadata (headers and messages).
- **CLI Commands**:
  - `init`: Initialize a new repository.
  - `cat-file`: Provide content or type and size information for repository objects.
  - `hash-object`: Compute object ID and optionally create a blob from a file.

## Getting Started

### Prerequisites

- Go 1.25.3 or later.

### Installation

Clone the repository and build the project:

```bash
git clone https://github.com/Blue-Onion/pygo.git
cd pygo
```

### Usage

You can run the CLI using `go run cmd/main.go` or using the provided `Makefile`.

#### Initialize a Repository

```bash
go run cmd/main.go init [path]
```

#### Hash a File

```bash
go run cmd/main.go hash-object [-t type] <file>
```

#### Inspect an Object

```bash
go run cmd/main.go cat-file <type> <object_sha>
```

## Project Structure

- `cmd/`: CLI entry point and command implementations.
- `hanlder/`: Core logic for Git objects and repository management.
  - `object/`: Object serialization, deserialization, and hashing.
  - `repo/`: Repository creation and lookup logic.
- `main.go`: Test script for the `Commit` object.
- `Makefile`: Convenient shortcuts for running and testing.

## License

MIT
