package errors

// BusinessError ...
type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New returns a new BusinessError
func New(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// Error returns error message
func (be *BusinessError) Error() string {
	return be.Message
}

// Wrap wraps error
func (be *BusinessError) Wrap(err error) *BusinessError {
	return &BusinessError{
		Code:    be.Code,
		Message: be.Message + ": " + err.Error(),
	}
}

// InvalidJSON means invalid json
var InvalidJSON = New(400000, "invalid json")

// FailedToGetToken means failed to get token
var FailedToGetToken = New(400100, "Failed to get token")

// FailedToGetOAuth2Provider means failed to get oauth2 provider
var FailedToGetOAuth2Provider = New(400101, "Failed to get oauth2 provider")

// InvalidCaptcha means invalid captcha
var InvalidCaptcha = New(400200, "验证码无效")

// UserLoginFailed means user login failed
var UserLoginFailed = New(400201, "登录失败")

// FailedToGetUser means failed to get user
var FailedToGetUser = New(400202, "Failed to get user")

// FailedToGetMenus means failed to get menus
var FailedToGetMenus = New(400203, "获取菜单失败")

// FailedToGetPermissions means failed to get permissions
var FailedToGetPermissions = New(400206, "获取权限失败")

// FailedToGetApps means failed to get apps
var FailedToGetApps = New(400204, "获取应用失败")

// FailedToGetUsers means failed to get users
var FailedToGetUsers = New(400205, "获取用户失败")
