package cli

import "os"

func GetEnv(key string, def string) (val string) {
	val, ok := os.LookupEnv(key)
	if !ok { val = def }
	return
}

func IsEnvSet(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}