package config

import (
	def "StakeBackendGoTest/pkg/const"
)

type Adapter struct {
	AdapterType     def.AdapterType `json:"adapter" yml:"adapter"`
	MockMSGCount    int             `json:"mock_msg_count" yml:"mock_msg_count"`
	MSGIntervalSecs int             `json:"msg_interval_secs" yml:"msg_interval_secs"`
}

type Config struct {
	Adapter *Adapter    `json:"adapter" yml:"adapter"`
	Addr    string      `json:"addr" yml:"addr"`
	EnvMode def.EnvMode `json:"env_mode" yml:"env_mode"`
	// TODO: put more fields here, say, for db settings, etc.
}
