package client

import (
	"github.com/alienantfarm/anthive/utils/structs"
)

type job struct {
	*endpoint
}

type Job struct {
	structs.Job
	endpoint *job
}
