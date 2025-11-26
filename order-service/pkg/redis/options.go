package redis

import "time"

type Option func(*Redis)

func ConnAttempts(a int) Option {
	return func(r *Redis) {
		r.connAttempts = a
	}
}

func Timeout(t time.Duration) Option {
	return func(r *Redis) {
		r.connTimeout = t
	}
}

func Prefix(p string) Option {
	return func(r *Redis) {
		r.keyPrefix = p
	}
}
