package main

import (
	"context"
	"log"
	"os"

	anytypemcp "github.com/elct9620/anytype-mcp"
	"github.com/elct9620/anytype-mcp/pkg/anytype"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	anytype := anytype.New(os.Getenv("ANYTYPE_API_KEY"))
	anytypeMcp := anytypemcp.New(anytype)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "anytype",
		Title:   "Anytype MCP",
		Version: "v0.1.0", // x-release-please-version
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
