package client

import (
	"bytes"
	"encoding/json"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/utils"
	"net/http"
	"strconv"
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
	structs.Antling
	endpoint *antling
}

func NewAntling(id int, client *Client) *Antling {
	return &Antling{structs.Antling{id, nil}, client.Antling}
}

func (a *Antling) GetJobs() ([]*structs.Job, error) {
	resp, err := a.endpoint.Client.Get(utils.Urlize(a.endpoint.Url, strconv.Itoa(a.Id)))
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusOK}
	}
	a.Jobs = []*structs.Job{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(a)
	return a.Jobs, err
}

func (a *Antling) Update() error {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(a)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"PATCH", utils.Urlize(a.endpoint.Url, strconv.Itoa(a.Id)), buf,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.endpoint.Client.Do(req)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusOK}
	}
	return err
}
