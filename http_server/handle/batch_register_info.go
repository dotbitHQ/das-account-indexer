package handle

import (
	"github.com/dotbitHQ/das-lib/common"
	code "github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqBatchRegisterInfo struct {
	BatchAccount []string `json:"batch_account" binding:"required,max=50"`
}

type RespBatchRegisterInfo struct {
	List []BatchRegisterInfoRecord `json:"list"`
}

type BatchRegisterInfoRecord struct {
	Account     string `json:"account"`
	CanRegister bool   `json:"can_register"`
}

func (h *HttpHandle) BatchRegisterInfo(ctx *gin.Context) {
	var (
		funcName = "BatchRegisterInfo"
		req      ReqBatchRegisterInfo
		apiResp  code.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doBatchRegisterInfo(&req, &apiResp); err != nil {
		log.Error("doBatchRegisterInfo err:", err.Error(), funcName)
	}
	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doBatchRegisterInfo(req *ReqBatchRegisterInfo, apiResp *code.ApiResp) error {
	accIds := make([]string, 0, len(req.BatchAccount))
	for _, v := range req.BatchAccount {
		accIds = append(accIds, common.Bytes2Hex(common.GetAccountIdByAccount(v)))
	}
	existAccIds, err := h.DbDao.GetAccountIdsByAccIds(accIds)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, err.Error())
		return err
	}
	existAccIdsMap := make(map[string]struct{}, len(existAccIds))
	for _, v := range existAccIds {
		existAccIdsMap[v] = struct{}{}
	}

	list := make([]BatchRegisterInfoRecord, 0, len(req.BatchAccount))
	for idx, v := range accIds {
		record := BatchRegisterInfoRecord{
			Account: req.BatchAccount[idx],
		}
		if _, ok := existAccIdsMap[v]; ok {
			accountName := strings.ToLower(strings.TrimSuffix(req.BatchAccount[idx], common.DasAccountSuffix))
			accountName = common.Bytes2Hex(common.Blake2b([]byte(accountName))[:20])
			_, reserved := h.MapReservedAccounts[accountName]
			_, unavailable := h.MapUnAvailableAccounts[accountName]
			record.CanRegister = !reserved && !unavailable
		}
		list = append(list, record)
	}

	apiResp.ApiRespOK(RespBatchRegisterInfo{
		List: list,
	})
	return nil
}
