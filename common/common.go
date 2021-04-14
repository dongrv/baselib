package common

// 基于HTTP Status进行扩展状态码
// 业务场景一般复杂度较高，四位码可以满足扩展
const (
	Continue                     = 1000 // 1xxx
	SwitchingProtocols           = 1001
	StatusOK                     = 2000 // 2xxx
	Created                      = 2001
	Accepted                     = 2002
	NonAuthoritativeInformation  = 2003
	NoContent                    = 2004
	ResetContent                 = 2005
	PartialContent               = 2006
	ValidAuthToken               = 2007
	InvalidAuthToken             = 2008
	InvalidRequestFormat         = 2009
	InvalidArgument              = 2010
	CodeInvalid                  = 2011
	VerifySignFailed             = 2012
	RegisterUserFailed           = 2013
	IllegalUserFound			 = 2014
	MultipleChoices              = 3000 // 3xxx
	MovedPermanently             = 3001
	MovedTemporarily             = 3002
	SeeOther                     = 3003
	NotModified                  = 3004
	UseProxy                     = 3005
	BadRequest                   = 4000 // 4xxx
	Unauthorized                 = 4001
	PaymentRequired              = 4002
	Forbidden                    = 4003
	NotFound                     = 4004
	MethodNotAllowed             = 4005
	NotAcceptable                = 4006
	ProxyAuthenticationRequired  = 4007
	RequestTimeout               = 4008
	Conflict                     = 4009
	Gone                         = 4010
	LengthRequired               = 4011
	PreconditionFailed           = 4012
	RequestEntityTooLarge        = 4013
	RequestURITooLong            = 4014
	UnsupportedMediaType         = 4015
	RequestedRangeNotSatisfiable = 4016
	ExpectationFailed            = 4017
	TooManyConnections           = 4021
	UnprocessableEntity          = 4022
	Locked                       = 4023
	FailedDependency             = 4024
	UnorderedCollection          = 4025
	UpgradeRequired              = 4026
	CredentialsInvalid           = 4106
	TypeAssertionError           = 4107
	InternalServerError          = 5000 // 5xxx
	NotImplemented               = 5001
	BadGateway                   = 5002
	ServiceUnavailable           = 5003
	GatewayTimeout               = 5004
	HTTPVersionNotSupported      = 5005
)

var statusText = map[uint16]string{
	Continue:                     "Continue",                        // 1xxx	// 继续
	SwitchingProtocols:           "Switching_Protocols",             // 切换协议
	StatusOK:                     "Success",                         // 2xxx	// 成功
	Created:                      "Created",                         // 已创建
	Accepted:                     "Accepted",                        // 已接受
	NonAuthoritativeInformation:  "Non-Authoritative_Information",   // 非授权信息
	NoContent:                    "No_Content",                      // 无内容
	ResetContent:                 "Reset_Content",                   // 重置内容
	PartialContent:               "Partial_Content",                 // 处理部分内容
	ValidAuthToken:               "Valid_Authorization_Token",       // 有效的授权凭证
	InvalidAuthToken:             "Invalid_Authorization_Token",     // 无效的授权凭证
	InvalidRequestFormat:         "Invalid_Request_Format",          // 无效的请求格式
	InvalidArgument:              "INVALID_ARGUMENT",                // 无效参数
	CodeInvalid:                  "CODE_INVALID",                    // 验证码无效
	VerifySignFailed:             "Verify_Sign_Failed",              // 验证签名失败
	RegisterUserFailed:           "Register_User_Failed",            // 注册用户失败
	IllegalUserFound:			  "Illegal_User_Found",	             // 非法的用户
	MultipleChoices:              "Multiple_Choices",                // 3xxx	// 多种选择
	MovedPermanently:             "Moved_Permanently",               // 永久移动
	MovedTemporarily:             "Moved_Temporarily",               // 临时移动
	SeeOther:                     "See_Other",                       // 查看其他位置
	NotModified:                  "Not_Modified",                    // 未修改
	UseProxy:                     "Use_Proxy",                       // 使用代理
	BadRequest:                   "Bad_Request",                     // 4xxx	// 错误的请求
	Unauthorized:                 "UNAUTHORISED",                    // 未授权
	PaymentRequired:              "PAYMENT_REQUIRED",                // 预留状态
	Forbidden:                    "FORBIDDEN",                       // 服务器拒绝执行
	NotFound:                     "NOT_FOUND",                       // 在服务器找不到资源
	MethodNotAllowed:             "METHOD_NOT_ALLOWED",              // 指定的请求方法不能被用于请求相应的资源
	NotAcceptable:                "Not_Acceptable",                  // 不接受
	ProxyAuthenticationRequired:  "Proxy_Authentication_Required",   // 需要代理授权
	RequestTimeout:               "Request_Timeout",                 // 请求超时
	Conflict:                     "Conflict",                        // 冲突
	Gone:                         "Gone",                            // 资源永久删除
	LengthRequired:               "Length_Required",                 // 需要有效长度
	PreconditionFailed:           "Precondition_Failed",             // 未能满足条件
	RequestEntityTooLarge:        "Request_Entity_Too_Large",        // 实体数据大小超过服务器处理范围，拒绝处理
	RequestURITooLong:            "Request-URI_Too_Long",            // URI 过长
	UnsupportedMediaType:         "Unsupported_MediaType",           // 不支持的资源格式
	RequestedRangeNotSatisfiable: "Requested_Range_Not_Satisfiable", // 不满足请求范围
	ExpectationFailed:            "Expectation_Failed",              // 请求头 Expect 中指定的预期内容无法被服务器满足
	TooManyConnections:           "Too_Many_Connections",            // 连接数超过服务器限制
	UnprocessableEntity:          "Unprocessable_Entity",            // 请求格式现语义错误
	Locked:                       "Locked",                          // 当前资源被锁定
	FailedDependency:             "Failed_Dependency",               // 由于之前的某个请求发生的错误，导致当前请求失败
	UnorderedCollection:          "Unordered_Collection",            // 无序集合
	UpgradeRequired:              "Upgrade_Required",                // 客户端应当切换到TLS/1.0
	CredentialsInvalid:           "CREDENTIAL_INVALID",              // 凭证无效
	TypeAssertionError:           "Type_Assertion_Error",            // 类型断言错误
	InternalServerError:          "Internal_Server_Error",           // 5xxx	// 服务器内部错误
	NotImplemented:               "Not_Implemented",                 // 服务器不具备完成请求的功能
	BadGateway:                   "Bad_Gateway",                     // 网关错误
	ServiceUnavailable:           "Service_Unavailable",             // 服务不可用
	GatewayTimeout:               "Gateway_Timeout",                 // 网关超时
	HTTPVersionNotSupported:      "HTTP_Version_Not_Supported",      // 服务器不支持请求中所用的 HTTP 协议版本

}

func StatusText(code uint16) string {
	if _, ok := statusText[code]; ok {
		return statusText[code]
	} else {
		return "Undefined_Type"
	}
}

var (
	// 常用字符库
	NewLine = []byte{'\n'}
	Space   = []byte{' '}
	Ping    = "Ping"
	Pong    = "Pong"

	// 话术库
	Logged                = "You are already logged in."
	LoginSuccessful       = "Welcome, you have successfully logged in."
	LoggedOutSuccessful   = "You have been logged out."
	BindAccountSuccessful = "Congratulations on successfully binding your account."
	// ...

)
