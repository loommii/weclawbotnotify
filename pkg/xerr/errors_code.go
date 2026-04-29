package xerr

// 预定义一些业务逻辑的错误
var (
	Success           = NewCodeError(0, "请求成功")
	RequestParamError = NewCodeError(100000, "请求参数错误")
	JwtError          = NewCodeError(100001, "Unauthorized")
)

// 第三方的错误 统一用 20XXXX
var ( // wechat API Error 2000XX

)
