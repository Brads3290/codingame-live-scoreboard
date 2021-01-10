package errors

import "strings"

type CompositeError struct {
	internalErrors []error
}

func (c CompositeError) Error() string {
	compArr := make([]string, 0)

	for _, v := range c.internalErrors {
		compArr = append(compArr, v.Error())
	}

	return strings.Join(compArr, "; ")
}
