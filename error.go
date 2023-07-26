package MQWP

type mqwpError struct {
	message string
}

func (err *mqwpError) Error() string {
	return err.message
}

func newError(errMsg string) error {
	return &mqwpError{
		message: errMsg,
	}
}
