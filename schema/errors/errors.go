package errors

import (
	baseerrors "errors"
)

func New(err string) error {
	return baseerrors.New(err)
}

func NewComposite(errs ...error) CompositeError {
	c := CompositeError{
		internalErrors: errs,
	}

	return c
}

func NewItemNotFound(err string) ItemNotFound {
	return ItemNotFound{
		error: New(err),
	}
}
