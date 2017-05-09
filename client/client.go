package client

import (
	"bytes"
	"encoding/json"
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

type anthive struct {
	*endpoint
	Antling *antling
}

type antling struct {
	*endpoint
}

func (e *endpoint) Post(body io.Reader) (*http.Response, error) {
	return e.Client.Post(e.Url, "application/json", body)
}

func (a *antling) Create() (antling *utils.Antling, err error) {
	antling = &utils.Antling{}
	resp, err := a.Post(bytes.NewReader(nil))
	if resp.StatusCode != http.StatusCreated {
		err = &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusCreated}
	}
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(antling)
	if err != nil {
		return
	}
	return
}

func NewClient() *anthive {
	client := &anthive{newEndpoint(viper.GetString("Anthive"), nil), nil}
	client.Antling = &antling{newEndpoint("antlings", client.endpoint)}
	return client
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
