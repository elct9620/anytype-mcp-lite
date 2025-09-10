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
	Data []struct {
		ID      string `json:"id"`
		SpaceId string `json:"space_id"`
		Name    string `json:"name"`
		Type    struct {
			ID   string `json:"id"`
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"data"`
	Pagination struct {
		Total  int `json:"total"`
		Offset int `json:"offset"`
	} `json:"pagination"`
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
