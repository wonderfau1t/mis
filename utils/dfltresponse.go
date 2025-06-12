package utils

const (
	StatusError = "Error"
	StatusOK    = "OK"
)

type Response[T any] struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Result       T      `json:"result,omitempty"`
}

func NewErrorResponse(errMsg string) Response[any] {
	return Response[any]{
		Status:       StatusError,
		ErrorMessage: errMsg,
	}
}

func NewSuccessResponse[T any](result T) Response[T] {
	return Response[T]{
		Status: StatusOK,
		Result: result,
	}
}
