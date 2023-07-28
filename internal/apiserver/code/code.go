package code

import (
	"fmt"
	"net/http"
	"sync"
)

const (
	Success                  = 200
	BadRequest               = 400
	Unauthorized             = 401
	Forbidden                = 403
	TokenRelatedUserDisabled = 403001
	NotFound                 = 404
	NotImplement             = 405
	Internal                 = 500
	UserLoginErr             = 406000 // 登录失败，用户名/密码错误
	UserResetPwErr           = 406001 // 重置密码失败
	ExistRelatedApp          = 406002
	OnlyLastAdmin            = 406003 // 超管最低人数限制
	DuplicateServiceName     = 406004 // 服务名称重复
	DuplicateChannelName     = 406005 // 渠道名称重复
	DuplicateAppName         = 406006
	NonExistentRequestID     = 406007 // 请求ID不存在
	ChannelNotExist          = 406008 // 渠道不存在
	DuplicateAccount         = 406009 // 账户重复
	UserDisabled             = 406010
	TooManyPoints            = 406011 // 时间间距过小
	CantEditYourself         = 406012 //禁止编辑自身信息
	GetAppErr                = 406013 // 获取应用失败
	DeleteAppErr             = 406014 // 删除应用失败
	StreamConfigErr          = 406015 // 服务配置信息错误
	DeleteChannelErr         = 406016 // 删除渠道失败
	DeleteLogsErr            = 406017 // 删除日志失败
	GetChannelErr            = 406018 // 获取渠道失败
	GetImageErr              = 406019 // 获取图片失败
	GetImageTypeErr          = 406020 // 获取图片类型失败
	GetServiceErr            = 407000 // 获取服务失败
	ServiceCreateErr         = 407001 // service 创建失败
	ChannelCreateErr         = 407003 // channel 创建失败
	AppCreateErr             = 407004 // app 创建失败
	DelServiceErr            = 407002 // 删除服务失败
	ConfigIconErr            = 407005 // 配置图标失败
	LogGetListErr            = 407010 // 获取日志列表失败
	SumRequestErr            = 408001 // 获取表格数据失败
	GetTrialListErr          = 408002 // 获取体验列表数据失败
	GetUpstreamErr           = 408003 // 获取上游配置失败
	CreateRouteErr           = 408004 // 创建 route 失败
	RouteConfigErr           = 408005
	// UserRegistered       = 10101 // 账户已注册
	// UserNotExist         = 10102
	// UserCreateErr        = 10103
	// UserInvalid          = 10105 // 用户未激活邮件
	// UserActivationErr    = 10107
	// UserAppVerifyErr     = 10108
	// UserDisabled         = 10109 // 用户被禁用
	// EmailActivateFailed  = 10130
	// EmailSendErr         = 10131
	// EmailTokenInvalid    = 10132
	// EmailTokenCreateErr  = 10133
	// EmailSendServerErr   = 10134 // 邮件发送服务失效
	// DuplicateData        = 20101
	// SignExpire           = 10301 // sign 签名过期失效
	// SignNonceDuplicate   = 10302 // sign 签名随机串已使用
	// SignAppIdNot         = 10303 // sign 签名 AppID 不存在
	// SignErr              = 10304 // sign 签名错误
	// SignModelErr         = 10305 // 找不到当前用户下的模型 id 以及版本
	// ModelNotExist        = 10401 // 模型不存在
	// ModelLabelNotExist   = 10402 // 模型标签不存在
	// ModelVersionNotExist = 10403 // 模型版本不存在
	// ModelCreateIng       = 10404 // 模型创建中
	// ModelTrainIng        = 10405 // 模型训练中
	// ModelDeployIng       = 10406 // 模型部署中
	// ModelErr             = 10407 // 模型错误
	// ModelLimit           = 10408 // 模型数量超过限制
	// ModelStatus          = 10409 // 模型状态错误
	// ModelStartFailed     = 10410 // 模型启动中
	// ModelStatusUpdate    = 10411 // 模型状态变化，不能修改删除
	// ModelStop            = 10412 // 模型已停止
	// ModelStarting        = 10413 // 模型启动中
	// ModelQuotaErr        = 10414 // 模型api配额不足
	// ImageErr             = 10420 // 图片错误
	// ImageApiErr          = 10421 // api 接口识别图片错误
	// ImageNotExist        = 10422 // 图片未上传
)

// 定义常用的状态码和，错误信息，可以根据需要自行添加
func init() {
	register(Success, 200, "success")
	register(BadRequest, 400, "bad request")
	register(Unauthorized, 401, "unauthorized")
	register(Forbidden, 403, "forbidden")
	register(TokenRelatedUserDisabled, 403, "该用户在管理平台被禁用")
	register(NotFound, 404, "not found")
	register(NotImplement, 404, "not implement")
	register(Internal, 500, "internal server error")
	register(UserLoginErr, 406, "错误的账号/密码")
	register(UserResetPwErr, 406, "重置密码失败")
	register(ExistRelatedApp, 406, "存在相关应用")
	register(OnlyLastAdmin, 406, "仅剩余最后一个超管")
	register(DuplicateServiceName, 406, "存在重名服务")
	register(DuplicateChannelName, 406, "存在重名渠道")
	register(DuplicateAppName, 406, "存在重名应用")
	register(NonExistentRequestID, 406, "请求ID不存在")
	register(ChannelNotExist, 406, "渠道不存在")
	register(DuplicateAccount, 406, "账户重复")
	register(UserDisabled, 406, "用户已停用")
	register(TooManyPoints, 406, "时间间隔过小")
	register(CantEditYourself, 406, "禁止编辑自身信息")
	register(ServiceCreateErr, 406, "服务创建失败")
	register(SumRequestErr, 406, "获取表格数据失败")
	register(GetTrialListErr, 406, "获取体验列表数据失败")
	register(GetUpstreamErr, 406, "获取上游配置失败")
	register(CreateRouteErr, 406, "创建 route 失败")
	register(RouteConfigErr, 406, "解析 route 配置失败")
	register(GetServiceErr, 406, "获取服务失败")
	register(DelServiceErr, 406, "删除服务失败")
	register(GetAppErr, 406, "获取应用失败")
	register(DeleteAppErr, 406, "删除应用失败")
	register(StreamConfigErr, 406, "服务配置信息错误")
	register(DeleteChannelErr, 406, "删除渠道失败")
	register(DeleteLogsErr, 406, "删除日志失败")
	register(GetChannelErr, 406, "获取渠道失败")
	register(GetImageErr, 406, "获取图片错误")
	register(GetImageTypeErr, 406, "获取图片类型错误")
	register(ChannelCreateErr, 406, "渠道创建失败")
	register(AppCreateErr, 406, "应用创建失败")
	register(ConfigIconErr, 406, "配置图标失败")
	register(LogGetListErr, 406, "获取日志列表失败")
}

var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

func register(code int, httpStatus int, message string) {
	coder := &ErrCode{
		C:          code,
		HttpStatus: httpStatus,
		Message:    message,
	}

	MustRegister(coder)
}

// MustRegister 注册错误码
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("ErrUnknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder 解析错误状态码
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			msg := v.Error() // 如果有错误信息，则返回错误信息
			if msg != "" {
				return ErrCode{
					C:          coder.Code(),
					HttpStatus: coder.HTTPStatus(),
					Message:    msg,
				}
			}
			return coder
		}
	}

	return ErrCode{
		C:          http.StatusBadRequest,
		HttpStatus: http.StatusBadRequest,
		Message:    err.Error(),
	}
}

// Coder 返回状态码接口定义
type Coder interface {
	Code() int
	HTTPStatus() int
	String() string
}

type ErrCode struct {
	C          int
	HttpStatus int
	Message    string
}

func (coder ErrCode) Code() int {
	return coder.C
}

func (coder ErrCode) String() string {
	return coder.Message
}

func (coder ErrCode) HTTPStatus() int {
	if coder.HttpStatus == 0 {
		return http.StatusInternalServerError
	}

	return coder.HttpStatus
}

type withCode struct {
	err  string
	code int
}

func (w *withCode) Error() string {
	if w.err == "" {
		return ""
	}
	return w.err
}

func WithCodeMsg(code int, msg ...any) error {
	var errMsg string
	switch len(msg) {
	case 0:
	case 1:
		errMsg = msg[0].(string)
	default:
		errMsg = fmt.Sprintf(msg[0].(string), msg[1:]...)
	}
	return &withCode{
		err:  errMsg,
		code: code,
	}
}
