# Architecture

This document provides an overview of project architecture and core aspects.

## Structure

To keep the project organized and simple, we use the following structure:

```
|- cmd/          # Application entry points
    |- main.go   # Initialize and start the application
|- pkg/
    |- anytype/  # The Go client library for Anytype
|-  anytype.go   # The adapter implementation for the Anytype client to MCP server
```

The project root is the adapter implementation for the Anytype client to MCP server. The adapter related code will put here.

## Anytype Client

The Anytype client is implemented in the `pkg/anytype` directory. We following the Anytype API documentation but only defined necessary methods and deserialize the required JSON fields.

## MCP Server

The MCP server is initialized and started in the `cmd/main.go` file. It imports the Anytype adapter to register as MCP tools.

```go
anytypeMcp := anytypemcplite.New(anytype) # anytypemcplite is the adapter package located in the project root

server := mcp.NewServer(...)

mcp.AddTool(server, &mcp.Tool{Name: "search", Description: "search objects in anytype"}, anytypeMcp.Search)
mcp.AddTool(server, &mcp.Tool{Name: "get-object", Description: "get an object from anytype"}, anytypeMcp.GetObject)
```
