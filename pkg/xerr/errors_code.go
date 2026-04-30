package xerr

var (
	Success           = NewCodeError(0, "请求成功")
	RequestParamError = NewCodeError(100000, "请求参数错误")
	JwtError          = NewCodeError(100001, "Unauthorized")
)

// 注册错误 1101XX
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

// 登录错误 1102XX
var (
	LoginParamEmpty    = NewCodeError(110200, "用户名和密码不能为空")
	LoginUserNotFound  = NewCodeError(110201, "用户名或密码错误")
	LoginPasswordWrong = NewCodeError(110202, "用户名或密码错误")
	LoginQueryFailed   = NewCodeError(110203, "查询用户失败")
	LoginTokenFailed   = NewCodeError(110204, "生成令牌失败")
)

// 刷新令牌错误 1103XX
var (
	RefreshTokenInvalid  = NewCodeError(110300, "刷新令牌无效")
	RefreshTokenRevoked  = NewCodeError(110301, "刷新令牌已失效")
	RefreshTokenExpired  = NewCodeError(110302, "刷新令牌已过期")
	RefreshTokenStoreErr = NewCodeError(110303, "存储刷新令牌失败")
)

// 密码管理错误 1104XX
var (
	PasswordTooWeak = NewCodeError(110400, "密码强度不足，至少8位且包含大小写字母和数字")
)

// 第三方错误 20XXXX
var ()
