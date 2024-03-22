package common

import (
	"net/http"
	"strconv"
)

func ParseQueryParamInt(r *http.Request, paramName string, defaultValue int) int {
	paramValue := r.URL.Query().Get(paramName)
	if paramValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(paramValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func ParseQueryParamUint(r *http.Request, paramName string, defaultValue uint) uint {
	paramValue := r.URL.Query().Get(paramName)
	if paramValue == "" {
		return defaultValue
	}

	uintValue, err := strconv.ParseUint(paramValue, 10, 0)
	if err != nil {
		return defaultValue
	}

	return uint(uintValue)
}
