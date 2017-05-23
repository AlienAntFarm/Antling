package client

import (
	"bytes"
	"encoding/json"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

type Antling struct {
	*structs.Antling
	*endpoint
}

func newAntling(parent *endpoint) *Antling {
	return &Antling{
		&structs.Antling{utils.Config.Id, []*structs.Job{}},
		newEndpoint("antlings", parent),
	}
}

func (a *Antling) Create() (err error) {
	resp, err := a.Post(bytes.NewReader(nil))
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusCreated {
		err = &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusCreated}
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(a)
	if err != nil {
		return
	}
	return
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
	err := json.NewEncoder(buf).Encode(a.Antling)
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
	if glog.V(2) {
		glog.Infof("%q", a)
	}
	resp, err := a.endpoint.Client.Do(req)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = &utils.UnexpectedStatusCode{http.StatusOK, resp.StatusCode}
	}
	return err
}
