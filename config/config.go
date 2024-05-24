package config

import "time"

const (
	RPC_URL              = "https://cloudflare-eth.com"
	POLLING_INTERVAL     = 13 * time.Second
	POLLING_ERROR_TOLERANCE = 1 * time.Second
)
