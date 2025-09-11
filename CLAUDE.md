# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Anytype MCP Lite is a lightweight Model Context Protocol server providing read-only access to Anytype workspace data. It's designed to be token-efficient for local LLMs by minimizing context size and expanding object data inline.

## Commands

### Build & Development
```bash
make build          # Build the application to ./dist/anytype-mcp
make clean          # Clean build artifacts
go mod tidy         # Update dependencies
gofmt -w .          # Format Go code (auto-applied via hooks)
```

## Architecture

The codebase follows a clean layered architecture:

### Core Components

1. **pkg/anytype/** - HTTP client library for Anytype API
   - `anytype.go`: Main client with Bearer auth, connects to `http://127.0.0.1:31009`
   - `object.go`: Object retrieval operations
   - `search.go`: Search operations  
   - `error.go`: Structured error handling

2. **server/** - MCP protocol adapter
   - `anytype.go`: Server initialization and tool registration
   - `search.go`: Implements search tool with pagination
   - `get_object.go`: Implements get-object tool
   - `property.go`: Property type definitions (only text/date formats exposed)

3. **cmd/main.go** - Application entry point
   - Initializes Anytype client with API key from environment
   - Creates MCP server and registers tools
   - Starts stdio transport

### Key Design Decisions

- **Token Efficiency**: Only implements `search` and `get-object` tools to minimize context
- **Expanded Data**: Returns full object content instead of references to reduce API calls
- **Property Filtering**: Only exposes text and date properties to reduce noise
- **Local API**: Assumes Anytype runs locally on port 31009
- **API Version**: Uses Anytype API version `2025-05-20`

## MCP Tools

### search
Searches for Anytype objects by text query with optional pagination.
- Returns simplified metadata: ID, SpaceId, Name, Type

### get-object
Retrieves full object by ID and SpaceId.
- Returns markdown content and filtered properties (text/date only)

## Release Process

Uses Release Please with GoReleaser for automated releases:
- Conventional commits trigger changelog generation
- Multi-platform builds: Linux, Windows, macOS (amd64/arm64)
- Artifacts published to GitHub releases