package anytype

import (
	"context"
	"net/url"
	"strconv"
)

type SearchParams struct {
	Offset int `json:"offset"`
}

type SearchBody struct {
	Query string `json:"query"`
}

type SearchInput struct {
	Params SearchParams `json:"params"`
	Body   SearchBody   `json:"body"`
}

type SearchOutput struct {
	Data       []Object   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func (a *Anytype) Search(ctx context.Context, input SearchInput) (*SearchOutput, error) {
	var output SearchOutput

	params := url.Values{}
	params.Add("offset", strconv.Itoa(input.Params.Offset))

	err := a.Post(ctx, "/v1/search?"+params.Encode(), input.Body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
