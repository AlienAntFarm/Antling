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

func (a *antling) Create() (antling *utils.Antling, err error) {
	antling = &utils.Antling{}
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
