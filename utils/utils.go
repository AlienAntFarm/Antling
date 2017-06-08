package utils

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"os"
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

func MarshalJSONb(i interface{}) (b []byte) {
	var err error
	if b, err = json.Marshal(i); err != nil {
		glog.Errorf("%s when marshalling %q", err, i)
	}
	return b
}

func MarshalJSON(i interface{}) string {
	return string(MarshalJSONb(i))
}

func RemoveOnFail(path string, err error) {
	if err == nil {
		return
	}
	glog.Infof("removing file %s due to %s", path, err)
	// if remove fails, log the error but don't alter err
	if err := os.RemoveAll(path); err != nil {
		glog.Errorf("%s", err)
	}
}
