package zh

/** --------------------------------------------- *
 * @filename   zhCN/
 * @author     jinycoo <caojingyin@jinycoo.com>
 * @datetime   2019-05-25 17:49
 * @version    1.0.0
 * @desc       .....
 * ---------------------------------------------- */

var ZH_CN = map[int]string{
	0:    "ok",
	-1:  "未知类型",
	-2:  "状态错误",
	-3:  "状态错误",
	-4:  "参数错误",
	-5:  "数据错误",
	-6:  "未找到请求数据",
	-7:  "应用程序不存在或已被封禁",
	-8:  "Access Key错误",
	-9:  "API校验密匙错",
	-10: "调用方对该Method没有权限",

	-20: "Token错误",   // Token malformed
	-21: "Token尚未激活", // token is not valid yet
	-22: "Token已过期",  // token is expired
	-23: "Token无效",   // Token invalid
	-24: "Token签名无效", // Token Signature invalid

	-31: "Token 包含多个无效片段",   // token contains an invalid number of segments
	-32: "Token中没有包含Bearer", // token string should not contain 'bearer '
	-33: "无效的Token令牌头",      // failed to decode token header
	-34: "无法解组令牌头",          // failed to unmarshal token header
	-35: "无法解码令牌声明",         // failed to decode token claims
	-36: "无法解组令牌声明",         // failed to unmarshal token claims
	-37: "签名方法（alg）不可用",     // signing method (alg) is unavailable.
	-38: "签名方法（alg）未指定",     // signing method (alg) is unspecified.
	-39: "签名方法（alg）未指定",     // signing method %v is invalid

	-101: "账号未登录",
	-102: "账号被封停",
	-103: "网络错误，需要重新登录",
	-105: "验证码错误",
	-106: "账号未激活",
	-107: "账号非正式会员或在适应期",
	-108: "应用不存在或者被封禁",
	-110: "未绑定手机",
	-111: "csrf 校验失败",
	-112: "系统升级中",
	-113: "账号尚未实名认证",
	-114: "请先绑定手机",
	-115: "请先完成实名认证",
	-116: "手机号错误",
	-117: "银行卡号错误",
	-118: "网址错误",
	-304: "没有改动",
	-307: "撞车跳转",
	-400: "请求错误",
	-401: "未认证",
	-403: "访问权限不足",
	-404: "访问东东不存在",
	-405: "不支持该方法",
	-409: "冲突",

	-500: "服务器错误",
	-501: "服务错误,请稍候重试",
	-503: "过载保护,服务暂不可用",
	-504: "服务调用超时",
	-509: "超出限制",

	-600: "Token错误",
	-601: "Token尚未激活",
	-602: "登录状态已过期",
	-603: "无效Token",
	-604: "无效凭证",
	-605: "凭证已过期",
	-606: "凭证已使用",

	-616: "上传文件不存在",
	-617: "上传文件太大",
	-625: "登录失败次数太多",
	-626: "用户不存在",
	-627: "密码错误",
	-628: "密码太弱",
	-629: "用户名或密码错误",
	-632: "操作对象数量限制",
	-643: "被锁定",
	-650: "用户等级太低",
	-652: "重复的用户",
	-658: "登录状态已过期",
	-659: "Token参数缺失",
	-662: "密码时间戳过期",
	-688: "地理区域限制",
	-689: "版权限制",

	-700: "账号无效",
	-701: "账号未登录",
	-702: "账号未激活",
	-703: "账号被封停",
	-704: "账号不存在",
	-705: "账号已存在",
	-706: "Ticket已使用",
	-707: "账户异常", // AccountAbnormal

	-870: "数据配置错误",
	-871: "数据源未配置或字段名错误",

	-1000: "网络连接错误",
	-1001: "网络连接超时",
	-1200: "被降级过滤的请求",
	-1201: "rpc服务的client都不可用",
	-1202: "rpc服务的client没有授权",
	-1203: "rpc服务配置错误",
	-2201: "rpc服务的server端不可用",
	-2202: "rpc服务的server端错误",

	7000: "模板初始化错误",
	7001: "模板内容设置有误",
}
