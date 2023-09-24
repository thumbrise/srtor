package util

import (
	"os"
	"strconv"
)

func EnvGetBool(key string) bool {
	val, exists := os.LookupEnv(key)
	if !exists {
		return false
	}

	result, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return result

}
