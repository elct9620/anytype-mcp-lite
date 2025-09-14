package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/elct9620/anytype-mcp-lite/pkg/anytype"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetObject_Success(t *testing.T) {
	tests := []struct {
		name           string
		params         GetObjectParams
		mockResponse   *anytype.GetObjectOutput
		expectedResult *GetObjectResult
	}{
		{
			name: "basic object with text properties",
			params: GetObjectParams{
				ObjectId: "obj123",
				SpaceId:  "space456",
			},
			mockResponse: &anytype.GetObjectOutput{
				Object: anytype.Object{
					ID:       "obj123",
					SpaceId:  "space456",
					Name:     "Test Object",
					Markdown: "# Test Content\n\nThis is a test object.",
					Type: anytype.ObjectType{
						ID:   "type1",
						Key:  "note",
						Name: "Note",
					},
					Properties: []anytype.Property{
						{
							ID:     "prop1",
							Key:    "description",
							Name:   "Description",
							Format: "text",
							Text:   "A test description",
						},
						{
							ID:     "prop2",
							Key:    "title",
							Name:   "Title",
							Format: "text",
							Text:   "Test Title",
						},
					},
				},
			},
			expectedResult: &GetObjectResult{
				ObjectId: "obj123",
				SpaceId:  "space456",
				Markdown: "# Test Content\n\nThis is a test object.",
				Properties: []Property{
					{
						Name:   "Description",
						Format: "text",
						Value:  "A test description",
					},
					{
						Name:   "Title",
						Format: "text",
						Value:  "Test Title",
					},
				},
			},
		},
		{
			name: "object with date properties only",
			params: GetObjectParams{
				ObjectId: "obj789",
				SpaceId:  "space101",
			},
			mockResponse: &anytype.GetObjectOutput{
				Object: anytype.Object{
					ID:       "obj789",
					SpaceId:  "space101",
					Name:     "Date Object",
					Markdown: "Object with dates",
					Properties: []anytype.Property{
						{
							ID:     "prop1",
							Key:    "created",
							Name:   "Created Date",
							Format: "date",
							Date:   "2024-01-15",
						},
						{
							ID:     "prop2",
							Key:    "modified",
							Name:   "Modified Date",
							Format: "date",
							Date:   "2024-03-20",
						},
					},
				},
			},
			expectedResult: &GetObjectResult{
				ObjectId: "obj789",
				SpaceId:  "space101",
				Markdown: "Object with dates",
				Properties: []Property{
					{
						Name:   "Created Date",
						Format: "date",
						Value:  "2024-01-15",
					},
					{
						Name:   "Modified Date",
						Format: "date",
						Value:  "2024-03-20",
					},
				},
			},
		},
		{
			name: "object with mixed properties - filtering test",
			params: GetObjectParams{
				ObjectId: "mixed123",
				SpaceId:  "spaceMixed",
			},
			mockResponse: &anytype.GetObjectOutput{
				Object: anytype.Object{
					ID:       "mixed123",
					SpaceId:  "spaceMixed",
					Name:     "Mixed Properties Object",
					Markdown: "## Mixed content",
					Properties: []anytype.Property{
						{
							ID:     "prop1",
							Key:    "title",
							Name:   "Title",
							Format: "text",
							Text:   "Mixed Title",
						},
						{
							ID:     "prop2",
							Key:    "number",
							Name:   "Number Field",
							Format: "number",
						},
						{
							ID:     "prop3",
							Key:    "created",
							Name:   "Created",
							Format: "date",
							Date:   "2024-02-10",
						},
						{
							ID:     "prop4",
							Key:    "checkbox",
							Name:   "Is Active",
							Format: "checkbox",
						},
						{
							ID:     "prop5",
							Key:    "url",
							Name:   "Website",
							Format: "url",
						},
						{
							ID:     "prop6",
							Key:    "email",
							Name:   "Email",
							Format: "email",
						},
					},
				},
			},
			expectedResult: &GetObjectResult{
				ObjectId: "mixed123",
				SpaceId:  "spaceMixed",
				Markdown: "## Mixed content",
				Properties: []Property{
					{
						Name:   "Title",
						Format: "text",
						Value:  "Mixed Title",
					},
					{
						Name:   "Created",
						Format: "date",
						Value:  "2024-02-10",
					},
				},
			},
		},
		{
			name: "object with empty properties",
			params: GetObjectParams{
				ObjectId: "empty456",
				SpaceId:  "spaceEmpty",
			},
			mockResponse: &anytype.GetObjectOutput{
				Object: anytype.Object{
					ID:         "empty456",
					SpaceId:    "spaceEmpty",
					Name:       "Empty Properties Object",
					Markdown:   "Object without properties",
					Properties: []anytype.Property{},
				},
			},
			expectedResult: &GetObjectResult{
				ObjectId:   "empty456",
				SpaceId:    "spaceEmpty",
				Markdown:   "Object without properties",
				Properties: []Property{},
			},
		},
		{
			name: "object with special characters",
			params: GetObjectParams{
				ObjectId: "special-id_123",
				SpaceId:  "space_with-dashes",
			},
			mockResponse: &anytype.GetObjectOutput{
				Object: anytype.Object{
					ID:       "special-id_123",
					SpaceId:  "space_with-dashes",
					Name:     "Object with 特殊字符 & symbols",
					Markdown: "# Content with émojis 🎉 and symbols\n\n> Quote with special chars: €, £, ¥",
					Properties: []anytype.Property{
						{
							Name:   "Special Description",
							Format: "text",
							Text:   "Text with 中文 and émojis 🚀",
						},
					},
				},
			},
			expectedResult: &GetObjectResult{
				ObjectId: "special-id_123",
				SpaceId:  "space_with-dashes",
				Markdown: "# Content with émojis 🎉 and symbols\n\n> Quote with special chars: €, £, ¥",
				Properties: []Property{
					{
						Name:   "Special Description",
						Format: "text",
						Value:  "Text with 中文 and émojis 🚀",
					},
				},
			},
		},
		{
			name: "object without markdown content",
			params: GetObjectParams{
				ObjectId: "nomd789",
				SpaceId:  "spaceNoMd",
			},
			mockResponse: &anytype.GetObjectOutput{
				Object: anytype.Object{
					ID:       "nomd789",
					SpaceId:  "spaceNoMd",
					Name:     "No Markdown Object",
					Markdown: "",
					Properties: []anytype.Property{
						{
							Name:   "Description",
							Format: "text",
							Text:   "Object without markdown",
						},
					},
				},
			},
			expectedResult: &GetObjectResult{
				ObjectId: "nomd789",
				SpaceId:  "spaceNoMd",
				Markdown: "",
				Properties: []Property{
					{
						Name:   "Description",
						Format: "text",
						Value:  "Object without markdown",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				expectedPath := "/v1/spaces/" + tt.params.SpaceId + "/objects/" + tt.params.ObjectId
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
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

			mcpResult, result, err := app.GetObject(context.Background(), req, tt.params)

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

func TestGetObject_ErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		params        GetObjectParams
		statusCode    int
		mockError     *anytype.Error
		expectedError string
	}{
		{
			name: "object not found",
			params: GetObjectParams{
				ObjectId: "nonexistent",
				SpaceId:  "space123",
			},
			statusCode: http.StatusNotFound,
			mockError: &anytype.Error{
				Code:    "not_found",
				Message: "Object not found",
				Object:  "object",
				Status:  404,
			},
			expectedError: "The object object returned an error: Object not found (code: not_found, status: 404)",
		},
		{
			name: "unauthorized access",
			params: GetObjectParams{
				ObjectId: "obj123",
				SpaceId:  "space456",
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
		{
			name: "bad request - missing object id",
			params: GetObjectParams{
				ObjectId: "",
				SpaceId:  "space123",
			},
			statusCode: http.StatusBadRequest,
			mockError: &anytype.Error{
				Code:    "invalid_request",
				Message: "ObjectId cannot be empty",
				Object:  "object",
				Status:  400,
			},
			expectedError: "The object object returned an error: ObjectId cannot be empty (code: invalid_request, status: 400)",
		},
		{
			name: "internal server error",
			params: GetObjectParams{
				ObjectId: "obj123",
				SpaceId:  "space456",
			},
			statusCode: http.StatusInternalServerError,
			mockError: &anytype.Error{
				Code:    "internal_error",
				Message: "Database connection failed",
				Object:  "object",
				Status:  500,
			},
			expectedError: "The object object returned an error: Database connection failed (code: internal_error, status: 500)",
		},
		{
			name: "forbidden access",
			params: GetObjectParams{
				ObjectId: "restricted123",
				SpaceId:  "privateSpace",
			},
			statusCode: http.StatusForbidden,
			mockError: &anytype.Error{
				Code:    "forbidden",
				Message: "Access denied to this object",
				Object:  "object",
				Status:  403,
			},
			expectedError: "The object object returned an error: Access denied to this object (code: forbidden, status: 403)",
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

			mcpResult, result, err := app.GetObject(context.Background(), req, tt.params)

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

func TestGetObject_PropertyFiltering(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &anytype.GetObjectOutput{
			Object: anytype.Object{
				ID:      "test123",
				SpaceId: "space456",
				Properties: []anytype.Property{
					{Name: "Text1", Format: "text", Text: "value1"},
					{Name: "Date1", Format: "date", Date: "2024-01-01"},
					{Name: "Number1", Format: "number"},
					{Name: "Checkbox1", Format: "checkbox"},
					{Name: "URL1", Format: "url"},
					{Name: "Email1", Format: "email"},
					{Name: "Text2", Format: "text", Text: "value2"},
					{Name: "Date2", Format: "date", Date: "2024-02-02"},
					{Name: "Phone1", Format: "phone"},
					{Name: "Select1", Format: "select"},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	_, result, err := app.GetObject(context.Background(), req, GetObjectParams{
		ObjectId: "test123",
		SpaceId:  "space456",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCount := 4 // Both text and date properties
	if len(result.Properties) != expectedCount {
		t.Errorf("expected %d properties, got %d", expectedCount, len(result.Properties))
	}

	for _, prop := range result.Properties {
		if prop.Format != "text" && prop.Format != "date" {
			t.Errorf("unexpected property format %s, expected only text or date", prop.Format)
		}
	}
}

func TestGetObject_NetworkError(t *testing.T) {
	client := anytype.New("test-api-key", anytype.WithApiServer("http://localhost:99999"))
	app := New(client)
	req := &mcp.CallToolRequest{}

	mcpResult, result, err := app.GetObject(context.Background(), req, GetObjectParams{
		ObjectId: "obj123",
		SpaceId:  "space456",
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

func TestGetObject_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		default:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&anytype.GetObjectOutput{})
		}
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mcpResult, result, err := app.GetObject(ctx, req, GetObjectParams{
		ObjectId: "obj123",
		SpaceId:  "space456",
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

func TestGetObject_RequestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept header 'application/json', got %s", r.Header.Get("Accept"))
		}

		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("expected Authorization header 'Bearer test-api-key', got %s", r.Header.Get("Authorization"))
		}

		if r.Header.Get("Anytype-Version") != anytype.APIVersion {
			t.Errorf("expected Anytype-Version header %s, got %s", anytype.APIVersion, r.Header.Get("Anytype-Version"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&anytype.GetObjectOutput{
			Object: anytype.Object{
				ID:      "obj123",
				SpaceId: "space456",
			},
		})
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	_, _, err := app.GetObject(context.Background(), req, GetObjectParams{
		ObjectId: "obj123",
		SpaceId:  "space456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetObject_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&anytype.GetObjectOutput{})
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	_, result, err := app.GetObject(context.Background(), req, GetObjectParams{
		ObjectId: "obj123",
		SpaceId:  "space456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := &GetObjectResult{
		ObjectId:   "",
		SpaceId:    "",
		Markdown:   "",
		Properties: []Property{},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestGetObject_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid": json}`))
	}))
	defer server.Close()

	client := anytype.New("test-api-key", anytype.WithApiServer(server.URL))
	app := New(client)
	req := &mcp.CallToolRequest{}

	mcpResult, result, err := app.GetObject(context.Background(), req, GetObjectParams{
		ObjectId: "obj123",
		SpaceId:  "space456",
	})

	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}

	if mcpResult == nil || !mcpResult.IsError {
		t.Fatal("expected MCP error result")
	}

	if result != nil {
		t.Error("expected nil result on error")
	}
}
