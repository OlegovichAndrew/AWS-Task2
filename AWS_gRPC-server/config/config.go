package config

import "os"

func GetStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}

var GRPC_PORT = GetStringParameter("GRPC_PORT", ":9050")
