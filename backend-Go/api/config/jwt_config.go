package config

import "time"

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}
