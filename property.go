package anytypemcp

type Property struct {
	Name   string `json:"name" jsonschema:"the name of the property"`
	Format string `json:"format" jsonschema:"the format of the property"`
	Value  string `json:"value" jsonschema:"the value of the property"`
}
