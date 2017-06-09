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
	Id   int
	Jobs map[int]*structs.Job
	*endpoint
}

func newAntling(parent *endpoint) *Antling {
	return &Antling{
		utils.Config.Id,
		map[int]*structs.Job{},
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
	antling := struct {
		Id int `json:"id"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(antling); err != nil {
		return
	}
	a.Id = antling.Id
	return
}

func (a *Antling) GetJobs() ([]*structs.Job, error) {
	resp, err := a.endpoint.Client.Get(utils.Urlize(a.endpoint.Url, strconv.Itoa(a.Id)))
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusOK}
	}
	antling := &structs.Antling{Jobs: []*structs.Job{}}
	return antling.Jobs, json.NewDecoder(resp.Body).Decode(antling)
}

func (a *Antling) Update() error {
	antling := structs.Antling{a.Id, structs.ListJobs(a.Jobs)}
	buf := bytes.NewBuffer(utils.MarshalJSONb(antling))

	req, err := http.NewRequest(
		"PATCH", utils.Urlize(a.endpoint.Url, strconv.Itoa(a.Id)), buf,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if glog.V(2) {
		glog.Infof(utils.MarshalJSON(antling))
	}
	resp, err := a.endpoint.Client.Do(req)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = &utils.UnexpectedStatusCode{http.StatusOK, resp.StatusCode}
	}
	return err
}
