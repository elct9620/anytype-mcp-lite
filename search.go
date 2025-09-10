package anytypemcp

import (
	"context"

	"github.com/elct9620/anytype-mcp-lite/pkg/anytype"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchParams struct {
	Query  string `json:"query" jsonschema:"the search query"`
	Offset int    `json:"offset" jsonschema:"the offset for pagination"`
}

type SearchItem struct {
	ID      string `json:"id" jsonschema:"the id of the object"`
	SpaceId string `json:"space_id" jsonschema:"the space id of the object"`
	Name    string `json:"name" jsonschema:"the name of the object"`
	Type    string `json:"type" jsonschema:"the type of the object"`
}

type SearchResult struct {
	Data       []SearchItem `json:"data" jsonschema:"the objects returned by the search"`
	Pagination Pagination   `json:"pagination" jsonschema:"the pagination info"`
}

func (a *App) Search(ctx context.Context, req *mcp.CallToolRequest, params SearchParams) (*mcp.CallToolResult, *SearchResult, error) {
	res, err := a.anytype.Search(ctx, anytype.SearchInput{
		Params: anytype.SearchParams{
			Offset: params.Offset,
		},
		Body: anytype.SearchBody{
			Query: params.Query,
		},
	})
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: err.Error()},
			},
		}, nil, err
	}

	items := make([]SearchItem, len(res.Data))
	for i, item := range res.Data {
		items[i] = SearchItem{
			ID:      item.ID,
			SpaceId: item.SpaceId,
			Name:    item.Name,
			Type:    item.Type.Name,
		}
	}

	return nil, &SearchResult{
		Data: items,
		Pagination: Pagination{
			Total:  res.Pagination.Total,
			Offset: res.Pagination.Offset,
		},
	}, nil
}
