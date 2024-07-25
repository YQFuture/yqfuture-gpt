package consts

const (
	Success                   = 20000 // 操作成功
	Fail                      = 40000 // 操作失败
	Unauthorized              = 40001 // 登录失效
	IncorrectCaptcha          = 40002 // 图像验证码错误
	IncorrectVerificationCode = 40003 // 手机验证码错误
	AutomaticLoginFailure     = 40004 // 自动登录失败
	PhoneIsRegistered         = 40005 // 手机号已注册
	PhoneIsNotRegistered      = 40006 // 手机号未注册
	PhoneIsNotBound           = 40007 // 手机号未绑定
	UserNotInOrg              = 40008 // 用户不在组织中
	OrgNameIsExist            = 40009 // 组织名称已存在
)
