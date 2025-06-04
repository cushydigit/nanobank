package config

import "time"

const (
	TTL_REFRESH_TOKEN               = time.Minute * 60 * 6
	TTL_ACCESS_TOKEN                = time.Minute * 15
	TTL_CONFIRMATION_TOKEN          = time.Minute * 5
	TIMEOUT_GATEWAY_READ            = time.Second * 10
	TIMEOUT_GATEWAY_WRITE           = time.Second * 15
	TIMEOUT_GATEWAY_IDLE            = time.Second * 60
	TIMEOUT_GATEWAY_READ_HEADER     = time.Second * 5
	INTERNAL_TIMEOUT_CLIENT         = time.Second * 10
	INTERNAL_CONNECTION_TIMEOUT     = time.Second * 5
	INTERNAL_KEEPALIVE_TIMEOUT      = time.Second * 30
	INTERNAL_IDLCONNECTION_MAXCOUNT = 100
	INTERNAL_IDLCONNECTION_TIMEOUT  = time.Second * 90
	INTERNAL_TLSHANDSHAKE_TIMEOUT   = time.Second * 10
)
