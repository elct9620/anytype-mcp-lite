package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Anytype struct {
	apiKey    string
	apiServer string
}

type AnytypeOption func(*Anytype)

func New(apiKey string, opts ...AnytypeOption) *Anytype {
	anytype := &Anytype{
		apiKey:    apiKey,
		apiServer: "http://127.0.0.1:31009",
	}

	for _, opt := range opts {
		opt(anytype)
	}

	return anytype
}

func WithApiServer(apiServer string) AnytypeOption {
	return func(a *Anytype) {
		a.apiServer = apiServer
	}
}

func (a *Anytype) Get(ctx context.Context, path string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.apiServer+path, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Anytype-Version", "2025-05-20")
	req.Header.Add("Authorization", "Bearer "+a.apiKey)
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp Error
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return err
		}

		return &errResp
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (a *Anytype) Post(ctx context.Context, path string, payload any, result any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.apiServer+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Anytype-Version", "2025-05-20")
	req.Header.Add("Authorization", "Bearer "+a.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp Error
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return err
		}
		return &errResp
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
