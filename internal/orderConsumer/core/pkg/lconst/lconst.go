package lconst

import "time"

const (
	RedisDBOrder     = 2
	OrderTimeout     = 10 * time.Minute
	TimeoutCheckTick = 1 * time.Minute
)
