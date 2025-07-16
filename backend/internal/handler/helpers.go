package handler

import (
	"net/url"
	"strconv"
)

func getQueryInt(query url.Values, key string, defaultValue int) int {
	str := query.Get(key)
	if str == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return value
}

func getQueryString(query url.Values, key string, defaultValue string) string {
	str := query.Get(key)
	if str == "" {
		return defaultValue
	}
	return str
}
