package interceptor

type APIResponseDTO[T interface{}] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func NewDefaultSuccessResponse[T interface{}](data T) *APIResponseDTO[T] {
	return &APIResponseDTO[T]{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
}

func NewDefaultErrorResponse[T interface{}](msg string) *APIResponseDTO[T] {
	return &APIResponseDTO[T]{
		Code: 1,
		Msg:  msg,
	}
}
