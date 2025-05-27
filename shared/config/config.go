package config

import "time"

const (
	TTL_REFRESH_TOKEN = time.Minute * 15
	TTL_ACCESS_TOKEN  = time.Minute * 1
)
