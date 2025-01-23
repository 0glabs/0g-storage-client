package api

// General errors
var (
	ErrNil        = NewBusinessError(0, "Success")
	ErrValidation = NewBusinessError(1, "Invalid parameter")
	ErrInternal   = NewBusinessError(2, "Internal server error")
)

type BusinessError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{code, message, nil}
}

func NewBusinessErrorWithData(code int, message string, data interface{}) *BusinessError {
	return &BusinessError{code, message, data}
}

func (err *BusinessError) Error() string {
	return err.Message
}

func (be *BusinessError) WithData(data interface{}) *BusinessError {
	return NewBusinessErrorWithData(be.Code, be.Message, data)
}
