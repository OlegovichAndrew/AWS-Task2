package utils

import (
	"aws-http-server/config"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func IncreaseFileValue(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(file)

	data := strings.TrimSpace(string(b))
	i, err := strconv.Atoi(data)
	if err != nil {
		return err
	}

	f, err := os.Create(config.FILE_NAME)
	if err != nil {
		return err
	}
	defer f.Close()

	write := strconv.Itoa(i + 1)
	_, err = f.WriteString(write)
	return err
}

func IsTextFile(filename string) bool {
	result := strings.SplitAfter(filename, ".")
	return result[len(result)-1] == "txt"
}
