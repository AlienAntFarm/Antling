package client

import (
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"io"
	"net/http"
	"time"
)

type endpoint struct {
	*http.Client
	Url    string
	Parent *endpoint
}

func (e *endpoint) Post(body io.Reader) (*http.Response, error) {
	return e.Client.Post(e.Url, "application/json", body)
}

func newEndpoint(fragment string, parent *endpoint) *endpoint {
	var client *http.Client

	if parent != nil {
		fragment = utils.Urlize(parent.Url, fragment)
		client = parent.Client
	} else {
		client = &http.Client{Timeout: time.Second * 10}
	}
	return &endpoint{client, fragment, parent}
}

type client struct {
	*endpoint
	Antling *Antling
}

var c *client

func NewClient() *client {
	if c == nil {
		c = &client{newEndpoint(utils.Config.Anthive, nil), nil}
		c.Antling = newAntling(c.endpoint)
	}
	glog.V(2).Infof("asking for a new client %q", c)
	return c
}
