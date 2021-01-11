package errors

type ItemNotFound struct {
	error
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case ItemNotFound:
		return true
	default:
		return false
	}
}
