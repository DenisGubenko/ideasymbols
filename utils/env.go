package utils

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GetEnvVariable(name string) (value string) {
	value, exists := os.LookupEnv(name)

	if !exists {
		formattedMessage := "env variable '%s' does not provided"
		logrus.Errorf(formattedMessage, name)
		os.Exit(2)
	}

	return
}

func GetIntEnvVariable(name string) int {
	stringValue := GetEnvVariable(name)
	convertedValue, err := strconv.ParseInt(stringValue, 10, 32)

	if err != nil {
		formattedMessage := "can not convert %s env variable with value '%s' to integer"
		logrus.Errorf(formattedMessage, name, stringValue)
		os.Exit(2)
	}

	return int(convertedValue)
}
