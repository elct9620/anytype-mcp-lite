package anytype

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetObject_Success(t *testing.T) {
	tests := []struct {
		name         string
		input        GetObjectInput
		mockResponse GetObjectOutput
		expectedPath string
	}{
		{
			name: "basic object retrieval",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "obj123",
					SpaceId:  "space456",
				},
			},
			mockResponse: GetObjectOutput{
				Object: Object{
					ID:      "obj123",
					SpaceId: "space456",
					Name:    "Test Object",
					Type: ObjectType{
						ID:   "type1",
						Key:  "note",
						Name: "Note",
					},
				},
			},
			expectedPath: "/v1/spaces/space456/objects/obj123",
		},
		{
			name: "object with full properties",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "obj789",
					SpaceId:  "space101",
				},
			},
			mockResponse: GetObjectOutput{
				Object: Object{
					ID:       "obj789",
					SpaceId:  "space101",
					Name:     "Complex Object",
					Markdown: "# Title\n\nThis is markdown content.",
					Type: ObjectType{
						ID:   "type2",
						Key:  "page",
						Name: "Page",
					},
					Properties: []Property{
						{
							ID:     "prop1",
							Key:    "description",
							Name:   "Description",
							Format: "text",
							Text:   "Sample description",
						},
						{
							ID:     "prop2",
							Key:    "created_date",
							Name:   "Created Date",
							Format: "date",
							Date:   "2023-01-15",
						},
					},
				},
			},
			expectedPath: "/v1/spaces/space101/objects/obj789",
		},
		{
			name: "object with special characters in IDs",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "obj-with-dashes_123",
					SpaceId:  "space_with_underscores",
				},
			},
			mockResponse: GetObjectOutput{
				Object: Object{
					ID:      "obj-with-dashes_123",
					SpaceId: "space_with_underscores",
					Name:    "Object with Special Chars",
					Type: ObjectType{
						ID:   "type3",
						Key:  "task",
						Name: "Task",
					},
				},
			},
			expectedPath: "/v1/spaces/space_with_underscores/objects/obj-with-dashes_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET method, got %s", r.Method)
				}

				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(tt.mockResponse); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer server.Close()

			client := New("test-api-key", WithApiServer(server.URL))
			result, err := client.GetObject(context.Background(), tt.input)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(result, &tt.mockResponse) {
				t.Errorf("expected %+v, got %+v", tt.mockResponse, *result)
			}
		})
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

		if r.Header.Get("Anytype-Version") != APIVersion {
			t.Errorf("expected Anytype-Version header %s, got %s", APIVersion, r.Header.Get("Anytype-Version"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetObjectOutput{})
	}))
	defer server.Close()

	client := New("test-api-key", WithApiServer(server.URL))
	_, err := client.GetObject(context.Background(), GetObjectInput{
		Params: GetObjectParams{
			ObjectId: "obj123",
			SpaceId:  "space456",
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetObject_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       GetObjectInput
		statusCode  int
		errorResp   Error
		expectedErr string
	}{
		{
			name: "not found error",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "nonexistent",
					SpaceId:  "space123",
				},
			},
			statusCode: http.StatusNotFound,
			errorResp: Error{
				Code:    "not_found",
				Message: "Object not found",
				Object:  "object",
				Status:  404,
			},
			expectedErr: "The object object returned an error: Object not found (code: not_found, status: 404)",
		},
		{
			name: "unauthorized error",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "obj123",
					SpaceId:  "space456",
				},
			},
			statusCode: http.StatusUnauthorized,
			errorResp: Error{
				Code:    "unauthorized",
				Message: "Invalid API key",
				Object:  "auth",
				Status:  401,
			},
			expectedErr: "The object auth returned an error: Invalid API key (code: unauthorized, status: 401)",
		},
		{
			name: "bad request error",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "",
					SpaceId:  "space123",
				},
			},
			statusCode: http.StatusBadRequest,
			errorResp: Error{
				Code:    "invalid_request",
				Message: "ObjectId cannot be empty",
				Object:  "object",
				Status:  400,
			},
			expectedErr: "The object object returned an error: ObjectId cannot be empty (code: invalid_request, status: 400)",
		},
		{
			name: "internal server error",
			input: GetObjectInput{
				Params: GetObjectParams{
					ObjectId: "obj123",
					SpaceId:  "space456",
				},
			},
			statusCode: http.StatusInternalServerError,
			errorResp: Error{
				Code:    "internal_error",
				Message: "Database connection failed",
				Object:  "object",
				Status:  500,
			},
			expectedErr: "The object object returned an error: Database connection failed (code: internal_error, status: 500)",
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
			_, err := client.GetObject(context.Background(), tt.input)

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

func TestGetObject_NetworkError(t *testing.T) {
	client := New("test-api-key", WithApiServer("http://localhost:99999"))

	_, err := client.GetObject(context.Background(), GetObjectInput{
		Params: GetObjectParams{
			ObjectId: "obj123",
			SpaceId:  "space456",
		},
	})

	if err == nil {
		t.Fatal("expected network error, got nil")
	}

	if _, ok := err.(*Error); ok {
		t.Fatalf("expected network error, got API error: %v", err)
	}
}

func TestGetObject_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid": json}`))
	}))
	defer server.Close()

	client := New("test-api-key", WithApiServer(server.URL))
	_, err := client.GetObject(context.Background(), GetObjectInput{
		Params: GetObjectParams{
			ObjectId: "obj123",
			SpaceId:  "space456",
		},
	})

	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}

	if _, ok := err.(*Error); ok {
		t.Fatalf("expected JSON decode error, got API error: %v", err)
	}
}

func TestGetObject_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetObjectOutput{})
	}))
	defer server.Close()

	client := New("test-api-key", WithApiServer(server.URL))
	result, err := client.GetObject(context.Background(), GetObjectInput{
		Params: GetObjectParams{
			ObjectId: "obj123",
			SpaceId:  "space456",
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := &GetObjectOutput{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestGetObject_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(GetObjectOutput{})
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := New("test-api-key", WithApiServer(server.URL))
	_, err := client.GetObject(ctx, GetObjectInput{
		Params: GetObjectParams{
			ObjectId: "obj123",
			SpaceId:  "space456",
		},
	})

	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}

	if err == context.Canceled {
		return
	}

	if _, ok := err.(*Error); ok {
		t.Fatalf("expected network/context error, got API error: %v", err)
	}
}
