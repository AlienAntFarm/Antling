package client

import (
	"bytes"
	"encoding/json"
	"github.com/alienantfarm/antling/utils"
	"net/http"
)

type antling struct {
	*endpoint
}

func (a *antling) Create() (antling *Antling, err error) {
	antling = &Antling{endpoint: a}
	resp, err := a.Post(bytes.NewReader(nil))
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusCreated {
		err = &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusCreated}
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(antling)
	if err != nil {
		return
	}
	return
}

type Antling struct {
	Id       int    `json:"id"`
	Jobs     []*Job `json:"jobs"`
	endpoint *antling
}

func NewAntling(id int, client *Client) *Antling {
	return &Antling{id, nil, client.Antling}
}

func (a *Antling) GetJobs() ([]*Job, error) {
	resp, err := a.endpoint.Get(a.Id)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusOK}
	}
	a.Jobs = []*Job{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(a)
	return a.Jobs, err
}
