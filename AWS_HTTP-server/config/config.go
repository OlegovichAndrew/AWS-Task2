package config

import "os"

func GetStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}

var GRPC_ADDR = GetStringParameter("GRPC_ADDR", ":9050")
var HTTP_ADDR = GetStringParameter("HTTP_ADDR", ":8070")
var BUCKET_NAME = GetStringParameter("BUCKET_NAME", "ul.practice")
var FILE_NAME = GetStringParameter("FILE_NAME", "number.txt")
