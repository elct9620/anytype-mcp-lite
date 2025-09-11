# Testing

This document outlines the criteria for evaluating the quality of tests in Anytype MCP. For each test, we assert its quality passed over 80% rubric score.

## Criteria

Following are the criteria used to evaluate the quality of tests, review step by step and give reasoning to explain why implementation can get the score.

### Naming Conventions (1 point)

Use descriptive names for test functions that clearly indicate what is being tested.

- `TestFunctionName` for single scenario can cover all cases
- `TestFunctionName_Scenario` for multiple scenarios with different contexts

### Single Responsibility (1 point)

Each test function should focus on a single aspect of the functionality being tested. Avoid combining multiple assertions or scenarios in one test.

```go
func Test_Add_PositiveNumbers(t *testing.T) {
    tests := []struct {
        a, b     int
        expected int
    }{
        {1, 2, 3},
        {10, 20, 30},
    }

    for _, tt := range tests {
        t.Parallel()
        result := Add(tt.a, tt.b)
        if result != tt.expected {
            t.Errorf("expected %d, got %d", tt.expected, result)
        }
    }
}
```

- Do not create complex validation functions in tests
- The validation should happen in for loop or subtests

### Table-driven Tests (1 point)

To testing multiple scenarios with same logic, we use table-driven tests. This makes it easy to add new test cases and improves readability.

```go
func Test_Add(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 1, 2, 3},
        {"negative numbers", -1, -2, -3},
        {"mixed numbers", -1, 1, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("expected %d, got %d", tt.expected, result)
            }
        })
    }
}
```

- Create anonymous struct with private fields
- Use `t.Run` to run subtests for each case
- Use `t.Parallel()` to run tests in parallel for better performance

### Use if instead of assert (1 point)

We do not use `testify/assert` package. Instead, we use standard `if` statements to keep simplicity and avoid external dependencies.

```go
func Test_Add(t *testing.T) {
    result := Add(1, 2)
    if result != 3 {
        t.Errorf("expected 3, got %d", result)
    }
}
```

- Use `if` statements to check conditions
- Use deep equality checks for complex types with `reflect.DeepEqual`

### Mocking HTTP Servers (1 point)

We prefer to simulate real environment as much as possible. For HTTP clients, we use `httptest` package to create mock servers.

```go
func Test_FetchData(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"key":"value"}`))
    }))
    defer server.Close()

    client := &http.Client{}
    data, err := FetchData(client, server.URL)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    expected := map[string]string{"key": "value"}
    if !reflect.DeepEqual(data, expected) {
        t.Errorf("expected %v, got %v", expected, data)
    }
}
```

- Use `httptest` as primary tool for mocking HTTP servers
- For multiple scenarios, combine with table-driven tests

## Scoring

Each criterion only get the point when it is fully satisfied, otherwise get 0 point.
