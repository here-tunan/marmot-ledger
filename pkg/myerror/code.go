package myerror

type Code int

const (
	OK            Code = 200
	WrongParam    Code = 400
	Unauthorized  Code = 401
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
	case InternalError:
		return "服务器内部错误"
	default:
		return ""
	}
}
