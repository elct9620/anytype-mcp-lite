package anytypemcp

import "github.com/elct9620/anytype-mcp-lite/pkg/anytype"

type App struct {
	anytype *anytype.Anytype
}

func New(client *anytype.Anytype) *App {
	return &App{anytype: client}
}
