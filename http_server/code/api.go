package code

type ApiResp struct {
	ErrNo  Code        `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func (a *ApiResp) ApiRespErr(errNo Code, errMsg string) {
	a.ErrNo = errNo
	a.ErrMsg = errMsg
}

func (a *ApiResp) ApiRespOK(data interface{}) {
	a.ErrNo = ApiCodeSuccess
	a.Data = data
}

func ApiRespErr(errNo Code, errMsg string) ApiResp {
	return ApiResp{
		ErrNo:  errNo,
		ErrMsg: errMsg,
		Data:   nil,
	}
}
