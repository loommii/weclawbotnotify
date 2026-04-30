package xerr

// 预定义一些业务逻辑的错误
var (
	Success           = NewCodeError(0, "请求成功")
	RequestParamError = NewCodeError(100000, "请求参数错误")
	JwtError          = NewCodeError(100001, "Unauthorized")
)

// weclawbotnotify-api 服务错误 1101XX
var (
	RegisterParamEmpty    = NewCodeError(110100, "用户名和密码不能为空")
	RegisterUsernameExist = NewCodeError(110101, "用户名已存在")
	RegisterQueryFailed   = NewCodeError(110102, "查询用户失败")
	RegisterHashFailed    = NewCodeError(110103, "密码哈希失败")
	RegisterInsertFailed  = NewCodeError(110104, "创建用户失败")
	RegisterGetIdFailed   = NewCodeError(110105, "获取用户ID失败")
	RegisterTokenFailed   = NewCodeError(110106, "生成令牌失败")
	RegisterClosed        = NewCodeError(110107, "注册已关闭，V1版本仅支持单用户")
)

// 第三方的错误 统一用 20XXXX
var ( // wechat API Error 2000XX

)
