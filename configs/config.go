package config

import (
	def "StakeBackendGoTest/pkg/const"
)

type Adapter struct {
	AdapterType def.AdapterType `json:"adapter"`
}

type Config struct {
	Adapter *Adapter    `json:"adapter"`
	Addr    string      `json:"addr"`
	EnvMode def.EnvMode `json:"env_mode"`
	// TODO: put more fields here, say, for db settings, etc.
}
