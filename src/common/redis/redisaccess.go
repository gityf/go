package redis

import "strings"

func IsRedisError(errStr string) bool {
	return !strings.Contains(errStr, "redigo: nil returned")
}