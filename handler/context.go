package handler

import (
	"net/http"
	"strconv"
)

const (
	DEFAULT_LIMIT = 20
	MAXIMUM_LIMIT = 100
)

// ToLimit returns the Limit from the current request.
//
// If the limit is not present, it will be set to DEFAULT_LIMIT.  Limits are
// capped at MAXIMUM_LIMIT.
func ToLimit(r *http.Request) int {
	if len(r.FormValue("limit")) == 0 {
		return DEFAULT_LIMIT
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		return DEFAULT_LIMIT
	}

	if limit > MAXIMUM_LIMIT {
		return MAXIMUM_LIMIT
	}

	return limit
}

// ToOffset returns the Offset from current request
// query if offset doesn't present set default offset
// equal to 0
func ToOffset(r *http.Request) int {
	if len(r.FormValue("offset")) == 0 {
		return 0
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		return 0
	}

	return offset
}
