package errs

import (
	"fmt"

	"github.com/hedhyw/go-tarantool-structure/internal/consts"
)

// InvalidTagValueError .
type InvalidTagValueError struct {
	field string
}

// NewInvalidTagValueError creates new InvalidTagValueError
func NewInvalidTagValueError(field string) error {
	return InvalidTagValueError{field}
}

func (e InvalidTagValueError) Error() string {
	return fmt.Sprintf(
		"Invalid value of tag %s for field %s",
		consts.TagName,
		e.field,
	)
}
