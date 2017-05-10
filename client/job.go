package client

type job struct {
	*endpoint
}

type Job struct {
	Id       int `json:"id"`
	endpoint *job
}
