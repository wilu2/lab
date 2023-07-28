package errors

type WithCode struct {
	Err  string
	Code int
}

func (w *WithCode) Error() string {
	if w.Err == "" {
		return ""
	}
	return w.Err
}

func WithCodeMsg(code int, msg ...string) error {
	var msgStr string
	if len(msg) == 1 {
		msgStr = msg[0]
	}
	return &WithCode{
		Err:  msgStr,
		Code: code,
	}
}
