package pqerr

import (
	"errors"

	"github.com/lib/pq"
)

func IsUniqueViolatesError(err error) bool {
	const code = "23505"

	if err == nil {
		return false
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false
	}

	return pqErr.Code == code
}
