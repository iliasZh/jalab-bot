package error_code

import (
	"errors"

	"github.com/lib/pq"
)

const (
	SuccessfulCompletion      pq.ErrorCode = "00000"
	ForeignKeyViolation       pq.ErrorCode = "23503"
	UniqueConstraintViolation pq.ErrorCode = "23505"
)

func GetFrom(err error) pq.ErrorCode {
	var errPQ *pq.Error
	if errors.As(err, &errPQ) {
		return errPQ.Code
	}

	return SuccessfulCompletion
}
