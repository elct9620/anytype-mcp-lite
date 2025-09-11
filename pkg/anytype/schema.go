package anytype

// ObjectType represents the type of an Anytype object
type ObjectType struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Property represents a property of an Anytype object
type Property struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	Format string `json:"format"`
	Date   string `json:"date,omitempty"`
	Text   string `json:"text,omitempty"`
}

// Pagination represents pagination information for search results
type Pagination struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
}

// Object represents an Anytype object
type Object struct {
	ID         string     `json:"id"`
	SpaceId    string     `json:"space_id,omitempty"`
	Name       string     `json:"name"`
	Markdown   string     `json:"markdown,omitempty"`
	Type       ObjectType `json:"type"`
	Properties []Property `json:"properties,omitempty"`
}
