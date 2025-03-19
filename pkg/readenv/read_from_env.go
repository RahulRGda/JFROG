package readenv

import (
	"os"
	"strconv"
	"strings"
)

var GetEnvInttest = getEnvIntImpl
var GetEnvInt64test = getEnvInt64Impl

func GetEnvString(key string, fallback ...string) string {
	value := os.Getenv(key)
	if len(value) == 0 && len(fallback) > 0 {
		return fallback[0]
	}
	return value
}

func GetEnvStrings(key string, defaultVal ...string) []string {
	if v := os.Getenv(key); len(v) != 0 {
		if strings.Contains(v, ",") {
			return strings.Split(v, ",")
		} else if strings.Contains(v, ";") {
			return strings.Split(v, ";")
		}
	}
	return defaultVal
}

func GetEnvInt(key string, fallback ...int) int {
	value, _ := strconv.Atoi(os.Getenv(key))
	if len(fallback) > 0 && value == 0 {
		return fallback[0]
	}
	return value
}

func GetEnvInt64(key string, fallback ...int64) int64 {
	value, _ := strconv.Atoi(os.Getenv(key))
	if len(fallback) > 0 && value == 0 {
		return fallback[0]
	}
	return int64(value)
}

func GetEnvBool(key string, fallback ...bool) bool {
	value := os.Getenv(key)
	v, _ := strconv.ParseBool(value)
	if len(fallback) > 0 && len(value) == 0 {
		return fallback[0]
	}
	return v
}

func getEnvIntImpl(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvInt64Impl(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// Export your original functions for external use.
func GetEnvIntOriginal(key string, defaultValue int) int {
	return getEnvIntImpl(key, defaultValue)
}

func GetEnvInt64Original(key string, defaultValue int64) int64 {
	return getEnvInt64Impl(key, defaultValue)
}
