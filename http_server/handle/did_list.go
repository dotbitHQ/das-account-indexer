package handle

import (
	"context"
	"das-account-indexer/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqDidList struct {
	core.ChainTypeAddress
	Pagination
	DidType tables.DidCellStatus `json:"did_type"`
}

type RespDidList struct {
	Total int64     `json:"total"`
	List  []DidData `json:"did_list"`
}

type DidData struct {
	Outpoint      string               `json:"outpoint"`
	AccountId     string               `json:"account_id"`
	Account       string               `json:"account"`
	Args          string               `json:"args"`
	ExpiredAt     uint64               `json:"expired_at"`
	DidCellStatus tables.DidCellStatus `json:"did_cell_status"`
}

func (h *HttpHandle) DidList(ctx *gin.Context) {
	var (
		funcName = "DidList"
		clientIp = GetClientIp(ctx)
		req      ReqDidList
		apiResp  http_api.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, toolib.JsonString(req))

	if err = h.doDidList(ctx, &req, &apiResp); err != nil {
		log.Error(ctx, "doDidList err:", err.Error(), funcName, clientIp)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doDidList(ctx context.Context, req *ReqDidList, apiResp *http_api.ApiResp) error {
	var resp RespDidList
	data := make([]DidData, 0)

	addrHex, err := req.FormatChainTypeAddress(h.DasCore.NetType(), true)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return fmt.Errorf("FormatChainTypeAddress err: %s", err.Error())
	} else if addrHex.DasAlgorithmId != common.DasAlgorithmIdAnyLock {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return nil
	}

	args := common.Bytes2Hex(addrHex.ParsedAddress.Script.Args)
	res, err := h.DbDao.QueryDidCell(args, req.DidType, req.GetLimit(), req.GetOffset())
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search did cell list err")
		return fmt.Errorf("SearchAccountList err: %s", err.Error())
	}
	for _, v := range res {
		temp := DidData{
			Outpoint:      v.Outpoint,
			Account:       v.Account,
			AccountId:     v.AccountId,
			Args:          v.Args,
			ExpiredAt:     v.ExpiredAt,
			DidCellStatus: req.DidType,
		}
		data = append(data, temp)
	}
	resp.List = data

	total, err := h.DbDao.QueryDidCellTotal(args, req.DidType)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "search did cell total err")
		return fmt.Errorf("QueryDidCellTotal err: %s", err.Error())
	}
	resp.Total = total
	apiResp.ApiRespOK(resp)
	return nil
}
