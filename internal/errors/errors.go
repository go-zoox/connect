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

var InvalidJSON = New(400000, "invalid json")
var FailedToGetToken = New(400100, "Failed to get token")
var FailedToGetOAuth2Provider = New(400101, "Failed to get oauth2 provider")
var InvalidCaptcha = New(400200, "验证码无效")
var UserLoginFailed = New(400201, "登录失败")
var FailedToGetUser = New(400202, "Failed to get user")
var FailedToGetMenus = New(400203, "获取菜单失败")
var FailedToGetApps = New(400204, "获取应用失败")
