package anytype

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const APIVersion = "2025-05-20"

type Transport struct {
	apiKey     string
	apiVersion string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Add("Anytype-Version", t.apiVersion)
	req.Header.Add("Authorization", "Bearer "+t.apiKey)
	req.Header.Add("Accept", "application/json")
	if req.Method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}
	return http.DefaultTransport.RoundTrip(req)
}

type Anytype struct {
	apiServer  string
	httpClient *http.Client
}

type AnytypeOption func(*Anytype)

func New(apiKey string, opts ...AnytypeOption) *Anytype {
	anytype := &Anytype{
		apiServer: "http://127.0.0.1:31009",
		httpClient: &http.Client{
			Transport: &Transport{
				apiKey:     apiKey,
				apiVersion: APIVersion,
			},
		},
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

	resp, err := a.httpClient.Do(req)
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

	resp, err := a.httpClient.Do(req)
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
