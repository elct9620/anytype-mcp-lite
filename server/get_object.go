package server

import (
	"context"

	"github.com/elct9620/anytype-mcp-lite/pkg/anytype"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetObjectParams struct {
	ObjectId string `json:"objectId" jsonschema:"the id of the object to get"`
	SpaceId  string `json:"spaceId,omitempty" jsonschema:"the space id of the object to get"`
}

type GetObjectResult struct {
	ObjectId   string     `json:"objectId" jsonschema:"the id of the object"`
	SpaceId    string     `json:"spaceId,omitempty" jsonschema:"the space id of the object"`
	Markdown   string     `json:"markdown" jsonschema:"the markdown content of the object"`
	Properties []Property `json:"properties" jsonschema:"the properties of the object"`
}

func (a *App) GetObject(ctx context.Context, req *mcp.CallToolRequest, params GetObjectParams) (*mcp.CallToolResult, *GetObjectResult, error) {
	res, err := a.anytype.GetObject(ctx, anytype.GetObjectInput{
		Params: anytype.GetObjectParams{
			ObjectId: params.ObjectId,
			SpaceId:  params.SpaceId,
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

	props := make([]Property, 0)
	for _, prop := range res.Object.Properties {
		if prop.Format != "text" && prop.Format != "date" {
			continue
		}

		p := Property{
			Name:   prop.Name,
			Format: prop.Format,
			Value:  prop.Text,
		}
		if prop.Format == "date" {
			p.Value = prop.Date
		}

		props = append(props, p)

	}

	return nil, &GetObjectResult{
		ObjectId:   res.Object.ID,
		SpaceId:    res.Object.SpaceId,
		Markdown:   res.Object.Markdown,
		Properties: props,
	}, nil
}
