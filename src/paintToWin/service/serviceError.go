package service

type ServiceError struct {
	ServiceName   string
	OperationName string

	ErrorCode    string
	ErrorMessage string
}

func (s ServiceError) Error() string {
	return s.ErrorMessage
}
