package server

type Pagination struct {
	Total  int `json:"total" jsonschema:"the total number of results"`
	Offset int `json:"offset" jsonschema:"the offset for pagination"`
}
