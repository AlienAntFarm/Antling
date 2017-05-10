package client

type job struct {
	*endpoint
}

// >>> this as to be shared
const (
	JOB_NEW = iota
	JOB_PENDING
	JOB_FINISH
	JOB_ERROR
)

var JOB_STATES = [...]string{
	"NEW",
	"PENDING",
	"FINISH",
	"ERROR",
}

// <<< this as to be shared

type Job struct {
	Id       int `json:"id"`
	State    int `json:"state"`
	endpoint *job
}
