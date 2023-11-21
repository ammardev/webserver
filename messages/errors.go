package messages

type HttpError interface {
	error
	Status() int
}

type NotImplementedErr struct {
}

func (NotImplementedErr) Error() string {
	return "Not Implemented"
}

func (NotImplementedErr) Status() int {
	return 501
}

type HttpVersionNotSupportedErr struct {
}

func (HttpVersionNotSupportedErr) Error() string {
	return "HTTP Version Not Supported"
}

func (HttpVersionNotSupportedErr) Status() int {
	return 505
}

type NotFoundErr struct {
}

func (NotFoundErr) Error() string {
	return "Not Found"
}

func (NotFoundErr) Status() int {
	return 404
}
