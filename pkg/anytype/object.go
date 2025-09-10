package anytype

import "context"

type GetObjectParams struct {
	ObjectId string `json:"objectId"`
	SpaceId  string `json:"spaceId"`
}

type GetObjectInput struct {
	Params GetObjectParams `json:"params"`
}

type GetObjectOutput struct {
	Object struct {
		ID       string `json:"id"`
		SpaceId  string `json:"space_id,omitempty"`
		Name     string `json:"name"`
		Markdown string `json:"markdown"`
		Type     struct {
			ID   string `json:"id"`
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"type"`
		Properties []struct {
			ID     string `json:"id"`
			Key    string `json:"key"`
			Name   string `json:"name"`
			Format string `json:"format"`
			Date   string `json:"date,omitempty"`
			Text   string `json:"text,omitempty"`
		} `json:"properties"`
	} `json:"object"`
}

func (a *Anytype) GetObject(ctx context.Context, input GetObjectInput) (*GetObjectOutput, error) {
	var output GetObjectOutput

	err := a.Get(ctx, "/v1/spaces/"+input.Params.SpaceId+"/objects/"+input.Params.ObjectId, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
