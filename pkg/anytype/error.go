package anytype

import "fmt"

var _ error = &Error{}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"The object %s returned an error: %s (code: %s, status: %d)",
		e.Object,
		e.Message,
		e.Code,
		e.Status,
	)
}
