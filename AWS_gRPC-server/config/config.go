package config

import "os"

func GetStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}

var GRPC_ADDR = GetStringParameter("GRPC_ADDR", ":9000")
var HTTP_ADDR = GetStringParameter("HTTP_ADDR", ":8080")
