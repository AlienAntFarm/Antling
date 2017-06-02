package lxc

import (
	"gopkg.in/lxc/go-lxc.v2"
)

var LOG_LEVELS = [...]lxc.LogLevel{
	lxc.CRIT,
	lxc.ERROR,
	lxc.WARN,
	lxc.NOTICE,
	lxc.INFO,
	lxc.TRACE,
}
