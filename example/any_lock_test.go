package example

import (
	"das-account-indexer/http_server/handle"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/parnurzeal/gorequest"
	"github.com/scorpiotzh/toolib"
	"testing"
)

var (
	TestUrl = "https://test-indexer.d.id/v1"
)

func doReq(url string, req, data interface{}) error {
	var resp http_api.ApiResp
	resp.Data = &data

	_, _, errs := gorequest.New().Post(url).SendStruct(&req).EndStruct(&resp)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}
	if resp.ErrNo != http_api.ApiCodeSuccess {
		return fmt.Errorf("%d - %s", resp.ErrNo, resp.ErrMsg)
	}
	return nil
}

func TestBatchReverseRecordV2(t *testing.T) {
	req := handle.ReqBatchReverseRecordV2{
		BatchKeyInfo: []core.ChainTypeAddress{
			{
				Type: "blockchain",
				KeyInfo: core.KeyInfo{
					CoinType: common.CoinTypeCKB,
					Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgpzk3ntzys3nuwmvnar2lrs54l9pat6wy3qqcmu76w",
				},
			},
			//{
			//	Type: "blockchain",
			//	KeyInfo: core.KeyInfo{
			//		CoinType: common.CoinTypeBTC,
			//		Key:      "tb1qumrp5k2es0d0hy5z6044zr2305pyzc978qz0ju",
			//	},
			//},
			//{
			//	Type: "blockchain",
			//	KeyInfo: core.KeyInfo{
			//		CoinType: common.CoinTypeCKB,
			//		Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgytmmrfg7aczevlxngqnr28npj2849erjyqqhe2guh",
			//	},
			//},
			//{
			//	Type: "blockchain",
			//	KeyInfo: core.KeyInfo{
			//		CoinType: common.CoinTypeBTC,
			//		Key:      "tb1pzl9nkuavvt303hly08u3ug0v55yd3a8x86d8g5jsrllsaell8j5s8gzedg",
			//	},
			//},
		},
	}
	url := TestUrl + "/batch/reverse/record"
	var data handle.RespBatchReverseRecordV2
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
}

func TestReverseRecordV2(t *testing.T) {
	req := handle.ReqReverseRecordV2{
		ChainTypeAddress: core.ChainTypeAddress{
			Type: "blockchain",
			KeyInfo: core.KeyInfo{
				CoinType: common.CoinTypeBTC,
				Key:      "tb1pzl9nkuavvt303hly08u3ug0v55yd3a8x86d8g5jsrllsaell8j5s8gzedg",
				//Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgytmmrfg7aczevlxngqnr28npj2849erjyqqhe2guh",
				//Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjzk3ntzys3nuwmvnar2lrs54l9pat6wy3qq5glj65",
				//Key:      "0x15a33588908cF8Edb27D1AbE3852Bf287Abd3891",
			},
		},
	}
	url := TestUrl + "/reverse/record"
	var data handle.RespReverseRecordV2
	if err := doReq(url, req, &data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(toolib.JsonString(&data))
}
