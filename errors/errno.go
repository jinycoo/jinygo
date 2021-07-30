package errors

/** --------------------------------------------- *
 * @filename   errors/errno.go
 * @author     jinycoo <caojingyin@jinycoo.com>
 * @datetime   2019-05-16 15:50
 * @version    1.0.0
 * @desc       公共错误代码 0除外  0具有更多项
 * ---------------------------------------------- */

var (
	OK = add(0) // 正确

	// Common Error
	UnknownType   = add(-1) // 未知类型
	StatusErr     = add(-2) // 状态错误
	StatusInvalid = add(-3) // 状态无效
	ParamsErr     = add(-4) // 参数错误
	DataErr       = add(-5) // 数据错误
	NotFoundData  = add(-6) // 未找到请求数据

	AppKeyInvalid      = add(-7)  // 应用程序不存在或已被封禁
	AccessKeyErr       = add(-8)  // Access Key错误
	SignCheckErr       = add(-9)  // API校验密匙错
	MethodNoPermission = add(-10) // 调用方对该Method没有权限

	AuthTokenErr              = add(-20) // Token错误 - Token malformed
	AuthTokenNotValidYet      = add(-21) // Token尚未激活 - token is not valid yet
	AuthTokenExpired          = add(-22) // Token已过期 - token is expired
	AuthTokenInvalid          = add(-23) // Token无效 - Token invalid
	AuthTokenSignatureInvalid = add(-24) // Token签名无效 - Token Signature invalid

	InvalidKey = add(-3000) //   无效Key key is invalid
	InvalidKeyType = add(-3001) // 无效Key类型 key is of invalid type
	HashUnavailable = add(-3002) // 请求的哈希函数不可用 - the requested hash function is unavailable


	AuthJwtSegmentsNumInvalid = add(-31) // Token 包含多个无效片段 - token contains an invalid number of segments
	AuthJwtBearerErr          = add(-32) // Token中没有包含Bearer - token string should not contain 'bearer '
	AuthJwtHeaderInvalid      = add(-33) // 无效的Token令牌头 failed to decode token header
	AuthJwtHeaderUnmarshalErr = add(-34) // 无法解组令牌头 failed to unmarshal token header
	AuthJwtClaimsErr          = add(-35) // 无法解码令牌声明 failed to decode token claims
	AuthJwtClaimsUnmarshalErr = add(-36) // 无法解组令牌声明 failed to unmarshal token claims
	AuthJwtAlgUnverfiable     = add(-37) // 签名方法（alg）不可用  signing method (alg) is unavailable.
	AuthJwtAlgUnspecified     = add(-38) // 签名方法（alg）未指定。  signing method (alg) is unspecified.
	AuthJwtAlgInvalid         = add(-39) // 签名方法（alg）未指定。 // signing method %v is invalid

	NoLogin                 = add(-101) // 账号未登录
	UserDisabled            = add(-102) // 账号被封停
	LogInAgain              = add(-103) // 网络错误，需要重新登录
	CaptchaErr              = add(-105) // 验证码错误
	UserInactive            = add(-106) // 账号未激活
	UserNoMember            = add(-107) // 账号非正式会员或在适应期
	AppDenied               = add(-108) // 应用不存在或者被封禁
	MobileNoVerfiy          = add(-110) // 未绑定手机
	CsrfNotMatchErr         = add(-111) // csrf 校验失败
	ServiceUpdate           = add(-112) // 系统升级中
	UserIDCheckInvalid      = add(-113) // 账号尚未实名认证
	UserIDCheckInvalidPhone = add(-114) // 请先绑定手机
	UserIDCheckInvalidCard  = add(-115) // 请先完成实名认证
	MobileNoErr             = add(-116) // 手机号错误
	BankCardNoErr           = add(-117) // 银行卡号错误
	WebSiteURLErr           = add(-118) // 网址错误

	NotModified       = add(-304) // 没有改动
	TemporaryRedirect = add(-307) // 撞车跳转
	RequestErr        = add(-400) // 请求错误
	Unauthorized      = add(-401) // 未认证
	AccessDenied      = add(-403) // 访问权限不足
	NothingFound      = add(-404) // 请求不存在
	MethodNotAllowed  = add(-405) // 不支持该方法

	Conflict = add(-409) // 冲突

	ServerErr          = add(-500) // 服务器错误
	ServerErrAgain     = add(-501) // 服务错误,请稍候重试
	ServiceUnavailable = add(-503) // 过载保护,服务暂不可用
	Deadline           = add(-504) // 服务调用超时
	LimitExceed        = add(-509) // 超出限制

	TokenMalformed   = add(-600) // Token错误 令牌格式错误 Token is malformed
	TokenSigningErr      = add(-601) // 由于签名问题无法验证令牌 - Token could not be verified because of signing problems
	TokenSignatureInvalid = add(-602) // 签名验证失败 Signature validation failed
	TokenExpired   = add(-603) // Token已过期
	TokenInvalid   = add(-604) // 无效Token
	TokenNotValidYet = add(-605) // Token尚未激活

	TicketInvalid  = add(-610) // 无效凭证
	TicketExpired  = add(-605) // 凭证已过期
	TicketConsumed = add(-606) // 凭证已使用

	FileNotExists         = add(-616) // 上传文件不存在
	FileTooLarge          = add(-617) // 上传文件太大
	FailedTooManyTimes    = add(-625) // 登录失败次数太多
	UserNotExist          = add(-626) // 用户不存在
	PasswordErr           = add(-627) // 密码错误
	PasswordTooLeak       = add(-628) // 密码太弱
	UsernameOrPasswordErr = add(-629) // 用户名或密码错误
	TargetNumberLimit     = add(-632) // 操作对象数量限制
	TargetBlocked         = add(-643) // 被锁定
	UserLevelLow          = add(-650) // 用户等级太低
	UserDuplicate         = add(-652) // 重复的用户

	PasswordHashExpires = add(-662) // 密码时间戳过期
	AreaLimit           = add(-688) // 地理区域限制
	CopyrightLimit      = add(-689) // 版权限制

	// Account - 账号相关
	AccountInvalid   = add(-700) // 账号无效
	AccountNoLogin   = add(-701) // 账号未登录
	AccountInactive  = add(-702) // 账号未激活
	AccountDisabled  = add(-703) // 账号被封停
	AccountNotExist  = add(-704) // 账号不存在
	AccountDuplicate = add(-705) // 账号已存在
	AccountConsumed  = add(-706) // 账号已登录
	AccountAbnormal  = add(-707) // 账户异常

	AddFailure    = add(-800) // 添加失败
	UpdateFailure = add(-801) // 更新失败
	DeleteFailure = add(-802) // 删除失败

	DataSourceConfigErr           = add(-870) // 数据配置错误
	DataSourceConfigFieldNotFound = add(-871) // 数据源未配置或字段名错误

	LinkErr      = add(-1000) // 网络连接错误
	LinkTimeout  = add(-1001) // 网络连接超时
	Degrade      = add(-1200) // 被降级过滤的请求
	RPCConfigErr = add(-1203) // rpc服务配置错误
	RPCNoClient  = add(-1201) // rpc服务的client都不可用
	RPCNoAuth    = add(-1202) // rpc服务的client没有授权
	RPCNoServer  = add(-2201) // rpc服务的server端不可用
	RPCServeErr  = add(-2202) // rpc服务的server端错误

	InvalidCredential     = add(40001) // 不合法的调用凭证
	InvalidGrantType      = add(40002) // 不合法的 grant_type
	InvalidOpenId         = add(40003) // 不合法的 OpenID
	InvalidMediaType      = add(40004) // 不合法的媒体文件类型
	InvalidMediaId        = add(40007) // 不合法的 media_id
	InvalidMessageType    = add(40008) // 不合法的 message_type
	InvalidImageSize      = add(40009) // 不合法的图片大小
	InvalidVoiceSize      = add(40010) // 不合法的语音大小
	InvalidVideoSize      = add(40011) // 不合法的视频大小
	InvalidThumbSize      = add(40012) // 不合法的缩略图大小
	InvalidAppId          = add(40013) // 不合法的 AppID
	InvalidAccessToken    = add(40014) // 不合法的 access_token
	InvalidCode           = add(40029) // 不合法或已过期的 code
	InvalidRefreshToken   = add(40030) // 不合法的 refresh_token
	InvalidTemplateIdSize = add(40036) // 不合法的 template_id 长度
	InvalidTemplateId     = add(40037) // 不合法的 template_id
	InvalidUrlSize        = add(40039) // 不合法的 url 长度
	InvalidUrlDomain      = add(40048) // 不合法的 url 域名
	InvalidUrl            = add(40066) // 不合法的 url

	AccessTokenExpired     = add(42001) // access_token 超时
	RefreshTokenExpired    = add(42002) // refresh_token 超时
	CodeExpired            = add(42003) // code 超时
	RequireGETMethod       = add(43001) // 需要使用 GET 方法请求
	RequirePOSTMethod      = add(43002) // 需要使用 POST 方法请求
	RequireHttps           = add(43003) // 需要使用 HTTPS
	RequireSubscribe       = add(43004) // 需要订阅关系
	AccessTokenMissing     = add(41001) // 缺失 access_token 参数
	AppIdMissing           = add(41002) // 缺失 AppId 参数
	RefreshTokenMissing    = add(41003) // 缺失 refresh_token 参数
	ApiUnauthorized        = add(50001) // 接口未授权
	ApiParamSizeOutOfLimit = add(45008) // 参数长度超过限制
	ApiFReqOutOfLimit      = add(45009) // 接口调动频率超过限制
	ApiLimit               = add(45011) // 频率限制

	TemplateInitErr       = add(7000) // 模板初始化错误
	TemplateErr           = add(7001) // 模板设置错误
	TemplateDataSourceErr = add(7002) // 模板数据源设置错误
)
