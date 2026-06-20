package myerror

type Code int

const (
	OK            Code = 200
	WrongParam    Code = 400
	Unauthorized  Code = 401
	Forbidden     Code = 403
	NotFound      Code = 404
	InternalError Code = 500
)

func (code Code) String() string {
	switch code {
	case OK:
		return "OK"
	case WrongParam:
		return "参数错误"
	case Unauthorized:
		return "鉴权失败"
	case Forbidden:
		return "无权限操作"
	case NotFound:
		return "资源不存在"
	case InternalError:
		return "服务器内部错误"
	default:
		return ""
	}
}
