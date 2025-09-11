package anytype

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSearch_Success(t *testing.T) {
	tests := []struct {
		name           string
		input          SearchInput
		mockResponse   SearchOutput
		expectedParams string
	}{
		{
			name: "basic search",
			input: SearchInput{
				Params: SearchParams{Offset: 0},
				Body:   SearchBody{Query: "test query"},
			},
			mockResponse: SearchOutput{
				Data: []struct {
					ID      string `json:"id"`
					SpaceId string `json:"space_id"`
					Name    string `json:"name"`
					Type    struct {
						ID   string `json:"id"`
						Key  string `json:"key"`
						Name string `json:"name"`
					} `json:"type"`
				}{
					{
						ID:      "obj1",
						SpaceId: "space1",
						Name:    "Test Object",
						Type: struct {
							ID   string `json:"id"`
							Key  string `json:"key"`
							Name string `json:"name"`
						}{
							ID:   "type1",
							Key:  "note",
							Name: "Note",
						},
					},
				},
				Pagination: struct {
					Total  int `json:"total"`
					Offset int `json:"offset"`
				}{
					Total:  1,
					Offset: 0,
				},
			},
			expectedParams: "offset=0",
		},
		{
			name: "search with offset",
			input: SearchInput{
				Params: SearchParams{Offset: 10},
				Body:   SearchBody{Query: "another query"},
			},
			mockResponse: SearchOutput{
				Data: []struct {
					ID      string `json:"id"`
					SpaceId string `json:"space_id"`
					Name    string `json:"name"`
					Type    struct {
						ID   string `json:"id"`
						Key  string `json:"key"`
						Name string `json:"name"`
					} `json:"type"`
				}{
					{
						ID:      "obj2",
						SpaceId: "space2",
						Name:    "Another Object",
						Type: struct {
							ID   string `json:"id"`
							Key  string `json:"key"`
							Name string `json:"name"`
						}{
							ID:   "type2",
							Key:  "page",
							Name: "Page",
						},
					},
				},
				Pagination: struct {
					Total  int `json:"total"`
					Offset int `json:"offset"`
				}{
					Total:  20,
					Offset: 10,
				},
			},
			expectedParams: "offset=10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/v1/search" {
					t.Errorf("expected path /v1/search, got %s", r.URL.Path)
				}

				if r.URL.RawQuery != tt.expectedParams {
					t.Errorf("expected query params %s, got %s", tt.expectedParams, r.URL.RawQuery)
				}

				var body SearchBody
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Query != tt.input.Body.Query {
					t.Errorf("expected query %s, got %s", tt.input.Body.Query, body.Query)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(tt.mockResponse); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := New("test-api-key", WithApiServer(server.URL))
			result, err := client.Search(context.Background(), tt.input)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(result, &tt.mockResponse) {
				t.Errorf("expected %+v, got %+v", tt.mockResponse, *result)
			}
		})
	}
}

func TestSearch_Pagination(t *testing.T) {
	tests := []struct {
		name       string
		offset     int
		totalItems int
		dataCount  int
	}{
		{"first page", 0, 100, 10},
		{"second page", 10, 100, 10},
		{"last page partial", 90, 95, 5},
		{"empty result", 100, 50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := SearchOutput{
					Data: make([]struct {
						ID      string `json:"id"`
						SpaceId string `json:"space_id"`
						Name    string `json:"name"`
						Type    struct {
							ID   string `json:"id"`
							Key  string `json:"key"`
							Name string `json:"name"`
						} `json:"type"`
					}, tt.dataCount),
					Pagination: struct {
						Total  int `json:"total"`
						Offset int `json:"offset"`
					}{
						Total:  tt.totalItems,
						Offset: tt.offset,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			client := New("test-api-key", WithApiServer(server.URL))
			result, err := client.Search(context.Background(), SearchInput{
				Params: SearchParams{Offset: tt.offset},
				Body:   SearchBody{Query: "test"},
			})

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if result.Pagination.Offset != tt.offset {
				t.Errorf("expected offset %d, got %d", tt.offset, result.Pagination.Offset)
			}

			if result.Pagination.Total != tt.totalItems {
				t.Errorf("expected total %d, got %d", tt.totalItems, result.Pagination.Total)
			}

			if len(result.Data) != tt.dataCount {
				t.Errorf("expected %d items, got %d", tt.dataCount, len(result.Data))
			}
		})
	}
}

func TestSearch_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		errorResp   Error
		expectedErr string
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			errorResp: Error{
				Code:    "invalid_request",
				Message: "Invalid search query",
				Object:  "search",
				Status:  400,
			},
			expectedErr: "The object search returned an error: Invalid search query (code: invalid_request, status: 400)",
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			errorResp: Error{
				Code:    "internal_error",
				Message: "Database connection failed",
				Object:  "search",
				Status:  500,
			},
			expectedErr: "The object search returned an error: Database connection failed (code: internal_error, status: 500)",
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			errorResp: Error{
				Code:    "unauthorized",
				Message: "Invalid API key",
				Object:  "auth",
				Status:  401,
			},
			expectedErr: "The object auth returned an error: Invalid API key (code: unauthorized, status: 401)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.errorResp)
			}))
			defer server.Close()

			client := New("test-api-key", WithApiServer(server.URL))
			_, err := client.Search(context.Background(), SearchInput{
				Params: SearchParams{Offset: 0},
				Body:   SearchBody{Query: "test"},
			})

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			apiErr, ok := err.(*Error)
			if !ok {
				t.Fatalf("expected *Error type, got %T", err)
			}

			if apiErr.Code != tt.errorResp.Code {
				t.Errorf("expected error code %s, got %s", tt.errorResp.Code, apiErr.Code)
			}

			if apiErr.Message != tt.errorResp.Message {
				t.Errorf("expected error message %s, got %s", tt.errorResp.Message, apiErr.Message)
			}

			if apiErr.Status != tt.errorResp.Status {
				t.Errorf("expected error status %d, got %d", tt.errorResp.Status, apiErr.Status)
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error string %s, got %s", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestSearch_NetworkError(t *testing.T) {
	client := New("test-api-key", WithApiServer("http://localhost:99999"))

	_, err := client.Search(context.Background(), SearchInput{
		Params: SearchParams{Offset: 0},
		Body:   SearchBody{Query: "test"},
	})

	if err == nil {
		t.Fatal("expected network error, got nil")
	}

	if _, ok := err.(*Error); ok {
		t.Fatalf("expected network error, got API error: %v", err)
	}
}

func TestSearch_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		input      SearchInput
		wantQuery  string
		wantOffset string
	}{
		{
			name: "empty query",
			input: SearchInput{
				Params: SearchParams{Offset: 0},
				Body:   SearchBody{Query: ""},
			},
			wantQuery:  "",
			wantOffset: "0",
		},
		{
			name: "special characters in query",
			input: SearchInput{
				Params: SearchParams{Offset: 5},
				Body:   SearchBody{Query: "test & query with spaces"},
			},
			wantQuery:  "test & query with spaces",
			wantOffset: "5",
		},
		{
			name: "large offset",
			input: SearchInput{
				Params: SearchParams{Offset: 1000000},
				Body:   SearchBody{Query: "test"},
			},
			wantQuery:  "test",
			wantOffset: "1000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				offset := r.URL.Query().Get("offset")
				if offset != tt.wantOffset {
					t.Errorf("expected offset %s, got %s", tt.wantOffset, offset)
				}

				var body SearchBody
				json.NewDecoder(r.Body).Decode(&body)
				if body.Query != tt.wantQuery {
					t.Errorf("expected query %s, got %s", tt.wantQuery, body.Query)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(SearchOutput{})
			}))
			defer server.Close()

			client := New("test-api-key", WithApiServer(server.URL))
			_, err := client.Search(context.Background(), tt.input)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
