package errors

import (
	baseerrors "errors"
)

func New(s string) error {
	return baseerrors.New(s)
}

func NewComposite(errs ...error) CompositeError {
	c := CompositeError{
		internalErrors: errs,
	}

	return c
}
