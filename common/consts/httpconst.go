package consts

const (
	Success                   = 20000 // 操作成功
	Fail                      = 40000 // 操作失败
	Unauthorized              = 40001 // 未授权
	IncorrectCaptcha          = 40002 // 图像验证码错误
	IncorrectVerificationCode = 40003 // 手机验证码错误
	AutomaticLoginFailure     = 40004 // 自动登录失败
)
