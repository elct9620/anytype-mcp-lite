package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/elct9620/anytype-mcp-lite/pkg/anytype"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSearch_Success(t *testing.T) {
	tests := []struct {
		name           string
		params         SearchParams
		mockResponse   *anytype.SearchOutput
		expectedResult *SearchResult
	}{
		{
			name: "basic search with results",
			params: SearchParams{
				Query:  "test query",
				Offset: 0,
			},
			mockResponse: &anytype.SearchOutput{
				Data: []anytype.Object{
					{
						ID:      "obj1",
						SpaceId: "space1",
						Name:    "Test Object 1",
						Type: anytype.ObjectType{
							ID:   "type1",
							Key:  "note",
							Name: "Note",
						},
					},
					{
						ID:      "obj2",
						SpaceId: "space2",
						Name:    "Test Object 2",
						Type: anytype.ObjectType{
							ID:   "type2",
							Key:  "page",
							Name: "Page",
						},
					},
				},
				Pagination: anytype.Pagination{
					Total:  2,
					Offset: 0,
				},
			},
			expectedResult: &SearchResult{
				Data: []SearchItem{
					{
						ID:      "obj1",
						SpaceId: "space1",
						Name:    "Test Object 1",
						Type:    "Note",
					},
					{
						ID:      "obj2",
						SpaceId: "space2",
						Name:    "Test Object 2",
						Type:    "Page",
					},
				},
				Pagination: Pagination{
					Total:  2,
					Offset: 0,
				},
			},
		},
		{
			name: "search with pagination offset",
			params: SearchParams{
				Query:  "paginated query",
				Offset: 10,
			},
			mockResponse: &anytype.SearchOutput{
				Data: []anytype.Object{
					{
						ID:      "obj11",
						SpaceId: "space1",
						Name:    "Result 11",
						Type: anytype.ObjectType{
							ID:   "type1",
							Key:  "task",
							Name: "Task",
						},
					},
				},
				Pagination: anytype.Pagination{
					Total:  50,
					Offset: 10,
				},
			},
			expectedResult: &SearchResult{
				Data: []SearchItem{
					{
						ID:      "obj11",
						SpaceId: "space1",
						Name:    "Result 11",
						Type:    "Task",
					},
				},
				Pagination: Pagination{
					Total:  50,
					Offset: 10,
				},
			},
		},
		{
			name: "empty search results",
			params: SearchParams{
				Query:  "no results",
				Offset: 0,
			},
			mockResponse: &anytype.SearchOutput{
				Data: []anytype.Object{},
				Pagination: anytype.Pagination{
					Total:  0,
					Offset: 0,
				},
			},
			expectedResult: &SearchResult{
				Data: []SearchItem{},
				Pagination: Pagination{
					Total:  0,
					Offset: 0,
				},
			},
		},
		{
			name: "search with large offset",
			params: SearchParams{
				Query:  "test",
				Offset: 1000,
			},
			mockResponse: &anytype.SearchOutput{
				Data: []anytype.Object{},
				Pagination: anytype.Pagination{
					Total:  100,
					Offset: 1000,
				},
			},
			expectedResult: &SearchResult{
				Data: []SearchItem{},
				Pagination: Pagination{
					Total:  100,
					Offset: 1000,
				},
			},
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

				offset := r.URL.Query().Get("offset")
				expectedOffset := fmt.Sprintf("%d", tt.params.Offset)
				if offset != expectedOffset {
					t.Errorf("expected offset %s, got %s", expectedOffset, offset)
				}

				var body anytype.SearchBody
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Query != tt.params.Query {
					t.Errorf("expected query %s, got %s", tt.params.Query, body.Query)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(tt.mockResponse); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
			app := New(client)
			req := &mcp.CallToolRequest{}

			mcpResult, result, err := app.Search(context.Background(), req, tt.params)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if mcpResult != nil {
				t.Fatal("expected nil MCP result for success")
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("expected result %+v, got %+v", tt.expectedResult, result)
			}
		})
	}
}

func TestSearch_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		params        SearchParams
		statusCode    int
		mockError     *anytype.Error
		expectedError string
	}{
		{
			name: "bad request",
			params: SearchParams{
				Query:  "test",
				Offset: 0,
			},
			statusCode: http.StatusBadRequest,
			mockError: &anytype.Error{
				Code:    "invalid_request",
				Message: "Invalid search parameters",
				Object:  "search",
				Status:  400,
			},
			expectedError: "The object search returned an error: Invalid search parameters (code: invalid_request, status: 400)",
		},
		{
			name: "internal server error",
			params: SearchParams{
				Query:  "test",
				Offset: 0,
			},
			statusCode: http.StatusInternalServerError,
			mockError: &anytype.Error{
				Code:    "internal_error",
				Message: "Database connection failed",
				Object:  "search",
				Status:  500,
			},
			expectedError: "The object search returned an error: Database connection failed (code: internal_error, status: 500)",
		},
		{
			name: "unauthorized",
			params: SearchParams{
				Query:  "test",
				Offset: 0,
			},
			statusCode: http.StatusUnauthorized,
			mockError: &anytype.Error{
				Code:    "unauthorized",
				Message: "Invalid API key",
				Object:  "auth",
				Status:  401,
			},
			expectedError: "The object auth returned an error: Invalid API key (code: unauthorized, status: 401)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if err := json.NewEncoder(w).Encode(tt.mockError); err != nil {
					t.Fatalf("failed to encode error response: %v", err)
				}
			}))
			defer server.Close()

			client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
			app := New(client)
			req := &mcp.CallToolRequest{}

			mcpResult, result, err := app.Search(context.Background(), req, tt.params)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error %s, got %s", tt.expectedError, err.Error())
			}

			if mcpResult == nil {
				t.Fatal("expected MCP error result")
			}

			if !mcpResult.IsError {
				t.Fatal("expected IsError to be true")
			}

			if len(mcpResult.Content) == 0 {
				t.Fatal("expected error content in MCP result")
			}

			textContent, ok := mcpResult.Content[0].(*mcp.TextContent)
			if !ok {
				t.Fatal("expected TextContent in MCP result")
			}

			if textContent.Text != tt.expectedError {
				t.Errorf("expected MCP error text %s, got %s", tt.expectedError, textContent.Text)
			}

			if result != nil {
				t.Error("expected nil result on error")
			}
		})
	}
}

func TestSearch_DataTransformation(t *testing.T) {
	tests := []struct {
		name         string
		anytypeObj   anytype.Object
		expectedItem SearchItem
	}{
		{
			name: "complete object",
			anytypeObj: anytype.Object{
				ID:      "123",
				SpaceId: "space123",
				Name:    "Test Object",
				Type: anytype.ObjectType{
					ID:   "type1",
					Key:  "document",
					Name: "Document",
				},
			},
			expectedItem: SearchItem{
				ID:      "123",
				SpaceId: "space123",
				Name:    "Test Object",
				Type:    "Document",
			},
		},
		{
			name: "object with empty type name",
			anytypeObj: anytype.Object{
				ID:      "456",
				SpaceId: "space456",
				Name:    "Another Object",
				Type: anytype.ObjectType{
					ID:   "type2",
					Key:  "unknown",
					Name: "",
				},
			},
			expectedItem: SearchItem{
				ID:      "456",
				SpaceId: "space456",
				Name:    "Another Object",
				Type:    "",
			},
		},
		{
			name: "object with special characters",
			anytypeObj: anytype.Object{
				ID:      "789",
				SpaceId: "space789",
				Name:    "Object with 特殊字符 & symbols",
				Type: anytype.ObjectType{
					ID:   "type3",
					Key:  "custom",
					Name: "Custom Type",
				},
			},
			expectedItem: SearchItem{
				ID:      "789",
				SpaceId: "space789",
				Name:    "Object with 特殊字符 & symbols",
				Type:    "Custom Type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := &anytype.SearchOutput{
					Data: []anytype.Object{tt.anytypeObj},
					Pagination: anytype.Pagination{
						Total:  1,
						Offset: 0,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
			app := New(client)
			req := &mcp.CallToolRequest{}

			_, result, err := app.Search(context.Background(), req, SearchParams{
				Query:  "test",
				Offset: 0,
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result.Data) != 1 {
				t.Fatalf("expected 1 item, got %d", len(result.Data))
			}

			if !reflect.DeepEqual(result.Data[0], tt.expectedItem) {
				t.Errorf("expected item %+v, got %+v", tt.expectedItem, result.Data[0])
			}
		})
	}
}

func TestSearch_MultipleResults(t *testing.T) {
	objects := []anytype.Object{
		{
			ID:      "obj1",
			SpaceId: "space1",
			Name:    "First",
			Type:    anytype.ObjectType{Name: "Type1"},
		},
		{
			ID:      "obj2",
			SpaceId: "space2",
			Name:    "Second",
			Type:    anytype.ObjectType{Name: "Type2"},
		},
		{
			ID:      "obj3",
			SpaceId: "space3",
			Name:    "Third",
			Type:    anytype.ObjectType{Name: "Type3"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &anytype.SearchOutput{
			Data: objects,
			Pagination: anytype.Pagination{
				Total:  3,
				Offset: 0,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	_, result, err := app.Search(context.Background(), req, SearchParams{
		Query:  "multiple",
		Offset: 0,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result.Data))
	}

	for i, obj := range objects {
		if result.Data[i].ID != obj.ID {
			t.Errorf("item %d: expected ID %s, got %s", i, obj.ID, result.Data[i].ID)
		}
		if result.Data[i].SpaceId != obj.SpaceId {
			t.Errorf("item %d: expected SpaceId %s, got %s", i, obj.SpaceId, result.Data[i].SpaceId)
		}
		if result.Data[i].Name != obj.Name {
			t.Errorf("item %d: expected Name %s, got %s", i, obj.Name, result.Data[i].Name)
		}
		if result.Data[i].Type != obj.Type.Name {
			t.Errorf("item %d: expected Type %s, got %s", i, obj.Type.Name, result.Data[i].Type)
		}
	}
}

func TestSearch_NetworkError(t *testing.T) {
	client := anytype.New("test-api-key", anytype.WithApiServer("http://localhost:99999"))
	app := New(client)
	req := &mcp.CallToolRequest{}

	mcpResult, result, err := app.Search(context.Background(), req, SearchParams{
		Query:  "test",
		Offset: 0,
	})

	if err == nil {
		t.Fatal("expected network error, got nil")
	}

	if mcpResult == nil {
		t.Fatal("expected MCP error result")
	}

	if !mcpResult.IsError {
		t.Fatal("expected IsError to be true")
	}

	if result != nil {
		t.Error("expected nil result on error")
	}
}

func TestSearch_RequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		params     SearchParams
		wantQuery  string
		wantOffset string
	}{
		{
			name: "empty query",
			params: SearchParams{
				Query:  "",
				Offset: 0,
			},
			wantQuery:  "",
			wantOffset: "0",
		},
		{
			name: "special characters in query",
			params: SearchParams{
				Query:  "test & query with spaces",
				Offset: 5,
			},
			wantQuery:  "test & query with spaces",
			wantOffset: "5",
		},
		{
			name: "large offset",
			params: SearchParams{
				Query:  "test",
				Offset: 1000000,
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

				var body anytype.SearchBody
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Query != tt.wantQuery {
					t.Errorf("expected query %s, got %s", tt.wantQuery, body.Query)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(&anytype.SearchOutput{
					Data:       []anytype.Object{},
					Pagination: anytype.Pagination{},
				})
			}))
			defer server.Close()

			client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
			app := New(client)
			req := &mcp.CallToolRequest{}

			_, _, err := app.Search(context.Background(), req, tt.params)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestSearch_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response that would be cancelled
		select {
		case <-r.Context().Done():
			return
		case <-r.Context().Done():
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&anytype.SearchOutput{})
		}
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mcpResult, result, err := app.Search(ctx, req, SearchParams{
		Query:  "test",
		Offset: 0,
	})

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if mcpResult == nil || !mcpResult.IsError {
		t.Fatal("expected MCP error result")
	}

	if result != nil {
		t.Error("expected nil result on error")
	}
}
