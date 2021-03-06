package en

/** --------------------------------------------- *
 * @filename   en/
 * @author     jinycoo <caojingyin@jinycoo.com>
 * @datetime   2019-05-25 18:19
 * @version    1.0.0
 * @desc       .....
 * ---------------------------------------------- */

var EN = map[int]string{
	0:    "ok",
	-1:  "Unknow type",
	-2:  "Status Error",
	-3:  "Status Error",
	-4:  "Params Error",
	-5:  "Data Error",
	-6:  "404 NOT FOUND",
	-7:  "App key invalid",
	-8:  "Access key error",
	-9:  "API sign check error",
	-10: "Method no permission",

	-20: "token malformed",
	-21: "token is not valid yet",
	-22: "token is expired",
	-23: "token invalid",
	-24: "token signature is invalid", // Token签名无效

	-31: "token contains an invalid number of segments",
	-32: "token string should not contain 'bearer '",
	-33: "failed to decode token header",
	-34: "failed to unmarshal token header",
	-35: "failed to decode token claims",
	-36: "failed to unmarshal token claims",
	-37: "signing method (alg) is unavailable",
	-38: "signing method (alg) is unspecified",
	-39: "signing method (alg) is invalid",

	-101: "NoLogin",                 // 账号未登录
	-102: "UserDisabled",            // 账号被封停
	-103: "LogInAgain",              // 网络错误，需要重新登录
	-105: "CaptchaErr",              // 验证码错误
	-106: "UserInactive",            // 账号未激活
	-107: "UserNoMember",            // 账号非正式会员或在适应期
	-108: "AppDenied",               // 应用不存在或者被封禁
	-110: "MobileNoVerfiy",          // 未绑定手机
	-111: "CsrfNotMatchErr",         // csrf 校验失败
	-112: "ServiceUpdate",           // 系统升级中
	-113: "UserIDCheckInvalid",      // 账号尚未实名认证
	-114: "UserIDCheckInvalidPhone", // 请先绑定手机
	-115: "UserIDCheckInvalidCard",  // 请先完成实名认证

	-304: "NotModified",       // 木有改动
	-307: "TemporaryRedirect", // 撞车跳转
	-400: "RequestErr",        // 请求错误
	-401: "Unauthorized",      // 未认证
	-403: "AccessDenied",      // 访问权限不足
	-404: "NothingFound",      // 啥都木有
	-405: "MethodNotAllowed",  // 不支持该方法
	-409: "Conflict",          // 冲突

	-500: "ServerErr",          // 服务器错误
	-503: "ServiceUnavailable", // 过载保护,服务暂不可用
	-504: "Deadline",           // 服务调用超时
	-509: "LimitExceed",        // 超出限制

	-616: "FileNotExists",         // 上传文件不存在
	-617: "FileTooLarge",          // 上传文件太大
	-625: "FailedTooManyTimes",    // 登录失败次数太多
	-626: "UserNotExist",          // 用户不存在
	-627: "Password Error",        // 密码错误
	-628: "PasswordTooLeak",       // 密码太弱
	-629: "UsernameOrPasswordErr", // 用户名或密码错误
	-632: "TargetNumberLimit",     // 操作对象数量限制
	-643: "TargetBlocked",         // 被锁定
	-650: "UserLevelLow",          // 用户等级太低
	-652: "UserDuplicate",         // 重复的用户
	-658: "AccessTokenExpires",    // Token 过期
	-659: "AccessTokenMissing",    // Token 参数缺失
	-662: "PasswordHashExpires",   // 密码时间戳过期
	-688: "AreaLimit",             // 地理区域限制
	-689: "CopyrightLimit",        // 版权限制

	-707: "AccountAbnormal", //
	-870: "DataSourceConfigErr",
	-871: "DataSourceConfigFieldNotFound",

	-1000: "LinkErr",      // 网络连接错误
	-1200: "Degrade",      // 被降级过滤的请求
	-1201: "RPCNoClient",  // rpc服务的client都不可用
	-1202: "RPCNoAuth",    // rpc服务的client没有授权
	-1203: "RPCConfigErr", // rpc服务配置错误
	-2201: "RPCNoServer",  // rpc服务的server端不可用
	-2202: "RPCServeErr",  // rpc服务的server端错误

	7000: "TemplateInitErr",    // 模板初始化错误
	7001: "TemplateSettingErr", // 模板设置错误
}
