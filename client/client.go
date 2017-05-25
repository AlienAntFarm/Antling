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

type Client struct {
	*endpoint
	Antling *Antling
	Images  *Images
}

var client *Client

func NewClient() *Client {
	if client == nil {
		client = &Client{newEndpoint(utils.Config.Anthive, nil), nil, nil}
		client.Antling = newAntling(client.endpoint)
		client.Images = newImages(client.endpoint)
	}
	glog.V(2).Infof("asking for a new client %q", client)
	return client
}
