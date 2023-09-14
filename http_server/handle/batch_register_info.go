package handle

import (
	"bytes"
	"das-account-indexer/config"
	"encoding/binary"
	"encoding/json"
	"github.com/dotbitHQ/das-lib/common"
	code "github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/minio/blake2b-simd"
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

func (h *HttpHandle) JsonRpcBatchRegisterInfo(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqBatchRegisterInfo
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doBatchRegisterInfo(&req[0], apiResp); err != nil {
		log.Error("doBatchReverseRecord err:", err.Error())
	}
}

func (h *HttpHandle) doBatchRegisterInfo(req *ReqBatchRegisterInfo, apiResp *code.ApiResp) error {
	accIds := make([]string, 0, len(req.BatchAccount))
	for _, v := range req.BatchAccount {
		accIds = append(accIds, common.Bytes2Hex(common.GetAccountIdByAccount(v)))
	}
	existAcc, err := h.DbDao.GetAccountByAccIds(accIds)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, err.Error())
		return err
	}
	existAccIdsMap := make(map[string]struct{}, len(existAcc))
	for _, v := range existAcc {
		if v.IsExpired() {
			continue
		}
		existAccIdsMap[v.AccountId] = struct{}{}
	}

	list := make([]BatchRegisterInfoRecord, 0, len(req.BatchAccount))
	for idx, v := range accIds {
		record := BatchRegisterInfoRecord{
			Account: req.BatchAccount[idx],
		}
		if _, ok := existAccIdsMap[v]; !ok {
			count := strings.Count(strings.TrimSuffix(record.Account, common.DasAccountSuffix), ".")
			accName := record.Account
			if !strings.HasSuffix(accName, common.DasAccountSuffix) {
				accName += common.DasAccountSuffix
			}
			switch count {
			case 0:
				// main account
				record.CanRegister, err = h.checkMainAccount(accName)
				if err != nil {
					apiResp.ApiRespErr(code.ApiCodeError500, err.Error())
					return err
				}
			case 1:
				// sub_account
				record.CanRegister, err = h.checkSubAccount(accName)
				if err != nil {
					apiResp.ApiRespErr(code.ApiCodeError500, err.Error())
					return err
				}
			}
		}
		list = append(list, record)
	}

	apiResp.ApiRespOK(RespBatchRegisterInfo{
		List: list,
	})
	return nil
}

func (h *HttpHandle) checkMainAccount(account string) (bool, error) {
	var accLen int
	_, accLen, err := common.GetDotBitAccountLength(account)
	if err != nil {
		return false, err
	}
	if accLen < 4 {
		return false, nil
	}
	accountName := strings.TrimSuffix(account, common.DasAccountSuffix)
	if strings.Contains(accountName, ".") {
		return false, nil
	}

	accCharStr, err := common.GetAccountCharSetList(account)
	if err != nil {
		return false, err
	}

	var accountCharStr string
	for _, v := range accCharStr {
		if v.Char == "" {
			return false, nil
		}
		switch v.CharSetName {
		case common.AccountCharTypeEmoji:
			if _, ok := common.CharSetTypeEmojiMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeDigit:
			if _, ok := common.CharSetTypeDigitMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeEn:
			if _, ok := common.CharSetTypeEnMap[v.Char]; v.Char != "." && !ok {
				return false, nil
			}
		case common.AccountCharTypeJa:
			if _, ok := common.CharSetTypeJaMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeRu:
			if _, ok := common.CharSetTypeRuMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeTr:
			if _, ok := common.CharSetTypeTrMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeVi:
			if _, ok := common.CharSetTypeViMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeTh:
			if _, ok := common.CharSetTypeThMap[v.Char]; !ok {
				return false, nil
			}
		case common.AccountCharTypeKo:
			if _, ok := common.CharSetTypeKoMap[v.Char]; !ok {
				return false, nil
			}
		default:
			return false, nil
		}
		accountCharStr += v.Char
	}

	if !strings.HasSuffix(accountCharStr, common.DasAccountSuffix) {
		accountCharStr += common.DasAccountSuffix
	}
	if !strings.EqualFold(account, accountCharStr) {
		return false, nil
	}

	if isDiff := common.CheckAccountCharTypeDiff(accCharStr); isDiff {
		return false, nil
	}

	accountName = strings.ToLower(accountName)
	accountName = common.Bytes2Hex(common.Blake2b([]byte(accountName))[:20])
	_, reserved := h.MapReservedAccounts[accountName]
	_, unavailable := h.MapUnAvailableAccounts[accountName]
	if reserved || unavailable {
		return !reserved && !unavailable, nil
	}

	if accLen >= config.Cfg.Das.AccountMinLength && accLen <= config.Cfg.Das.AccountMaxLength &&
		accLen >= config.Cfg.Das.OpenAccountMinLength && accLen <= config.Cfg.Das.OpenAccountMaxLength {

		tc, err := h.DasCore.GetTimeCell()
		if err != nil {
			return false, err
		}
		tcTimestamp := tc.Timestamp()
		openTimestamp := int64(1666094400)
		if config.Cfg.Server.Net != common.DasNetTypeMainNet {
			openTimestamp = 1665712800
		}
		// check dao char type
		isSameDaoCharType := true
		for i, v := range accCharStr {
			if v.Char == "." {
				break
			}
			if i == 0 {
				continue
			}
			if _, ok := OpenCharTypeMap[accCharStr[i].CharSetName]; !ok {
				isSameDaoCharType = false
				break
			}
			if accCharStr[i].CharSetName != accCharStr[i-1].CharSetName {
				isSameDaoCharType = false
				break
			}
		}
		if tcTimestamp >= openTimestamp && isSameDaoCharType {
			return true, nil
		}

		configRelease, err := h.DasCore.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsRelease)
		if err != nil {
			return false, err
		}
		luckyNumber, _ := configRelease.LuckyNumber()
		log.Info("config release lucky number: ", luckyNumber)
		if resNum, _ := Blake256AndFourBytesBigEndian([]byte(accountCharStr)); resNum <= luckyNumber {
			return true, nil
		}
	}
	return false, nil
}

func (h *HttpHandle) checkSubAccount(account string) (bool, error) {
	parentAccId := common.GetAccountIdByAccount(account)
	accInfo, err := h.DbDao.GetAccountInfoByAccountId(common.Bytes2Hex(parentAccId))
	if err != nil {
		return false, err
	}
	if accInfo.Id == 0 || accInfo.IsExpired() {
		return true, nil
	}
	return false, nil
}

var OpenCharTypeMap = map[common.AccountCharType]struct{}{
	common.AccountCharTypeEmoji: {},
	common.AccountCharTypeDigit: {},
	common.AccountCharTypeKo:    {},
	common.AccountCharTypeTh:    {},
	//common.AccountCharTypeTr:    {},
	//common.AccountCharTypeVi:    {},
}

func Blake256AndFourBytesBigEndian(data []byte) (uint32, error) {
	bys, err := Blake256(data)
	if err != nil {
		return 0, err
	}
	bytesBuffer := bytes.NewBuffer(bys[0:4])
	var res uint32
	if err = binary.Read(bytesBuffer, binary.BigEndian, &res); err != nil {
		return 0, err
	}
	return res, nil
}

func Blake256(data []byte) ([]byte, error) {
	tmpConfig := &blake2b.Config{
		Size:   32,
		Person: []byte("2021-07-22 12:00"),
	}
	hash, err := blake2b.New(tmpConfig)
	if err != nil {
		return nil, err
	}
	hash.Write(data)
	return hash.Sum(nil), nil
}
