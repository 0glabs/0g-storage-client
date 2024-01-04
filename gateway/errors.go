package gateway

// General errors
var (
	ErrNil            = newBusinessError(0, "ok")
	ErrValidation     = newBusinessError(1, "Invalid parameter")
	ErrInternalServer = newBusinessError(2, "Internal server error")
)

type BusinessError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewBusinessError(code int, message string, data interface{}) *BusinessError {
	return &BusinessError{code, message, data}
}

func newBusinessError(code int, message string) *BusinessError {
	return NewBusinessError(code, message, nil)
}

func (err *BusinessError) Error() string {
	return err.Message
}

func (be *BusinessError) WithData(data interface{}) *BusinessError {
	return NewBusinessError(be.Code, be.Message, data)
}
