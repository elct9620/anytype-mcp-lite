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
	Object Object `json:"object"`
}

func (a *Anytype) GetObject(ctx context.Context, input GetObjectInput) (*GetObjectOutput, error) {
	var output GetObjectOutput

	err := a.Get(ctx, "/v1/spaces/"+input.Params.SpaceId+"/objects/"+input.Params.ObjectId, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
