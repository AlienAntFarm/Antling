package client

import (
	"github.com/spf13/viper"
)

type anthive struct {
	*endpoint
	Antling *antling
}

func newAnthive() *anthive {
	a := &anthive{newEndpoint(viper.GetString("Anthive"), nil), nil}
	a.Antling = &antling{newEndpoint("antlings", a.endpoint)}
	return a
}
