package code

import "github.com/dotbitHQ/das-lib/http_api"

type ApiResp struct {
	ErrNo  http_api.ApiCode `json:"errno"`
	ErrMsg string           `json:"errmsg"`
	http_api.ApiResp
}

func (a *ApiResp) ApiRespErr(errNo http_api.ApiCode, errMsg string) {
	a.ErrNo = errNo
	a.ErrMsg = errMsg
	a.ApiResp.ErrNo = a.ErrNo
	a.ApiResp.ErrMsg = errMsg
}

func (a *ApiResp) ApiRespOK(data interface{}) {
	a.ErrNo = http_api.ApiCodeSuccess
	a.Data = data
	a.ApiResp.ErrNo = a.ErrNo
}
