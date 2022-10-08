package config

import (
	def "StakeBackendGoTest/pkg/const"
)

type Adapter struct {
	AdapterType def.AdapterType `json:"adapter" yml:"adapter"`
}

type Config struct {
	Adapter *Adapter    `json:"adapter" yml:"adapter"`
	Addr    string      `json:"addr" yml:"addr"`
	EnvMode def.EnvMode `json:"env_mode" yml:"env_mode"`
	// TODO: put more fields here, say, for db settings, etc.
}
