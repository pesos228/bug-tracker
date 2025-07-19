package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func decodeJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return false
	}
	return true
}

func encodeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to Encode DTO: %s", err.Error()), http.StatusInternalServerError)
	}
}
