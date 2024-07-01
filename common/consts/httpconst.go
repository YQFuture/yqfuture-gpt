package consts

const (
	Success                 = 200 // 操作成功
	BadRequest              = 400 // 错误的请求
	Unauthorized            = 401 // 未授权
	Forbidden               = 403 // 禁止访问
	NotFound                = 404 // 未找到
	MethodNotAllowed        = 405 // 请求方法不被允许
	NotAcceptable           = 406 // 不可接受的响应
	RequestTimeout          = 408 // 请求超时
	Conflict                = 409 // 冲突
	Gone                    = 410 // 资源已不存在
	InternalServerError     = 500 // 服务器内部错误
	NotImplemented          = 501 // 未实现的功能
	BadGateway              = 502 // 错误的网关
	ServiceUnavailable      = 503 // 服务不可用
	GatewayTimeout          = 504 // 网关超时
	HTTPVersionNotSupported = 505 // 不支持的HTTP版本
)
