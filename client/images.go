package client

import (
	"github.com/alienantfarm/antling/utils"
	"github.com/golang/glog"
	"io"
	"net/http"
)

type Images struct {
	*endpoint
}

func newImages(parent *endpoint) *Images {
	return &Images{
		newEndpoint(utils.IMAGES_PREFIX, parent),
	}
}

func (i *Images) Get(id string) (io.ReadCloser, error) {
	glog.Infof("retrieving image %s", id)
	resp, err := i.Client.Get(utils.Urlize(i.Url, id))
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, &utils.UnexpectedStatusCode{resp.StatusCode, http.StatusOK}
	}
	return resp.Body, err
}
