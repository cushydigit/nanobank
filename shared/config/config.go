package config

import "time"

const (
	TTL_REFRESH_TOKEN           = time.Minute * 60 * 6
	TTL_ACCESS_TOKEN            = time.Minute * 15
	TIMEOUT_GATEWAY_READ        = time.Second * 10
	TIMEOUT_GATEWAY_WRITE       = time.Second * 15
	TIMEOUT_GATEWAY_IDLE        = time.Second * 60
	TIMEOUT_GATEWAY_READ_HEADER = time.Second * 5
)
