package main

import (
	"context"
	"log"
	"os"

	"github.com/elct9620/anytype-mcp-lite/pkg/anytype"
	"github.com/elct9620/anytype-mcp-lite/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// x-release-please-start-version
const Version = "0.1.2"

// x-release-please-end

func main() {
	anytype := anytype.New(os.Getenv("ANYTYPE_API_KEY"))
	anytypeMcp := server.New(anytype)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "anytype",
		Title:   "Anytype MCP",
		Version: "v" + Version,
	},
		&mcp.ServerOptions{
			Instructions: "Provide read-only access to Anytype workspace. Help user to retrieve information from their Anytype.",
		},
	)
	mcp.AddTool(server, &mcp.Tool{Name: "search", Description: "search objects in anytype"}, anytypeMcp.Search)
	mcp.AddTool(server, &mcp.Tool{Name: "get-object", Description: "get an object from anytype"}, anytypeMcp.GetObject)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
