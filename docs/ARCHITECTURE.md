# Architecture

This document provides an overview of project architecture and core aspects.

## Structure

To keep the project organized and simple, we use the following structure:

```
|- cmd/          # Application entry points
    |- main.go   # Initialize and start the application
|- pkg/
    |- anytype/  # The Go client library for Anytype
|- server/       # The MCP server implementation
    |- anytype.go # The adapter implementation for the Anytype client to MCP server
```

## Anytype Client

The Anytype client is implemented in the `pkg/anytype` directory. We following the Anytype API documentation but only defined necessary methods and deserialize the required JSON fields.

## MCP Server

The MCP server is implemented in the `server` directory. It implements the MCP protocol to adapt the Anytype client to MCP tools.

## Command

The entrypoint of the application is in `cmd/main.go`. It initializes the Anytype client, creates an MCP server, and registers the Anytype adapter as a tool.

```go
anytype := anytype.New("your-anytype-api-key") # pkg/anytype
anytypeMcp := server.New(anytype) # server/anytype.go

server := mcp.NewServer(...)

mcp.AddTool(server, &mcp.Tool{Name: "search", Description: "search objects in anytype"}, anytypeMcp.Search)
mcp.AddTool(server, &mcp.Tool{Name: "get-object", Description: "get an object from anytype"}, anytypeMcp.GetObject)
```
