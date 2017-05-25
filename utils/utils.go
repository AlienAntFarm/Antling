package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Urlize(fragments ...string) string {
	return strings.Join(fragments, "/")
}

type UnexpectedStatusCode struct {
	Expected int
	Received int
}

func (usc *UnexpectedStatusCode) Error() string {
	return fmt.Sprintf(
		"unexpected status code, want: %d, got: %d", usc.Expected, usc.Received,
	)
}

func MarshalJSON(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}
