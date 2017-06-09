package utils

import (
	"github.com/alienantfarm/anthive/utils/structs"
)

type Job struct {
	*structs.Job
	Retries int
}

func (j *Job) SanitizeCwd() string {
	if j.Cwd != "" {
		return j.Cwd
	}
	if j.Image.Cwd != "" {
		return j.Image.Cwd
	}
	return "/"
}

func (j *Job) SanitizeEnv() []string {
	env := j.Image.Env
	if j.Image.Env == nil {
		env = []string{}
	}
	if j.Env != nil {
		env = append(env, j.Env...)
	}
	return env
}

func (j *Job) SanitizeCmd() []string {
	if j.Cmd != nil {
		return j.Cmd
	}
	if j.Image.Cmd != nil {
		return j.Image.Cmd
	}
	return []string{"cat", "/etc/hostname"}
}
