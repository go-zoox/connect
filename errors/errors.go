package errors

type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

func (be *BusinessError) Error() string {
	return be.Message
}

func (be *BusinessError) Wrap(err error) *BusinessError {
	return &BusinessError{
		Code:    be.Code,
		Message: be.Message + ": " + err.Error(),
	}
}

var InvalidCaptcha = New(400100, "验证码无效")
var UserLoginFailed = New(400101, "登录失败")
