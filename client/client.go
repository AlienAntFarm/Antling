package client

import (
	"github.com/alienantfarm/antling/utils"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"time"
)

type endpoint struct {
	*http.Client
	Url    string
	Parent *endpoint
}

var client *Client

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
	Antling *antling
}

func Get() *Client {
	if client == nil {
		client = &Client{newEndpoint(viper.GetString("Anthive"), nil), nil}
		client.Antling = &antling{newEndpoint("antlings", client.endpoint)}
	}

	return client
}
