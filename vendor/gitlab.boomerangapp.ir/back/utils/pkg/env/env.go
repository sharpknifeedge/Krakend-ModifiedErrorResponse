package env

import (
	"os"
	"strconv"
)

//Str get os enviernment with default value as string
func Str(key, defaultValue string) string {

	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

//Bool get os enviernment with default value as bool
func Bool(key string, defaultValue bool) bool {

	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val == "Y" || val == "YES" || val == "y" || val == "yes" || val == "1"
}

//Int get os enviernment with default value as int
func Int(key string, defaultValue int64) int64 {
	if os.Getenv(key) == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil {
		return defaultValue
	}

	return i
}

//Float get os enviernment with default value as float
func Float(key string, defaultValue float64) float64 {
	if os.Getenv(key) == "" {
		return defaultValue
	}

	i, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		return defaultValue
	}

	return i
}
