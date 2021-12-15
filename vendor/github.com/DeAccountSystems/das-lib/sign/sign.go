package sign

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/scorpiotzh/mylog"
)

var log = mylog.NewLogger("sign", mylog.LevelDebug)

type reqParam struct {
	Errno  int         `json:"errno"`
	Errmsg interface{} `json:"errmsg"`
	Data   interface{} `json:"data"`
}

type Client struct {
	ctx    context.Context
	client rpc.Client
}

func NewClient(ctx context.Context, apiUrl string) (*Client, error) {
	client, err := rpc.Dial(apiUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:    ctx,
		client: client,
	}, nil
}

func (c *Client) Client() rpc.Client {
	return c.client
}

func (c *Client) SignCkbMessage(ckbSignerAddress, message string) ([]byte, error) {
	if c.client != nil {
		if common.Has0xPrefix(message) {
			message = message[2:]
		}
		reply := reqParam{}
		param := struct {
			Address     string `json:"address"`
			CkbBuildRet string `json:"ckb_build_ret"`
			Tx          string `json:"tx"`
		}{
			Address:     ckbSignerAddress,
			CkbBuildRet: "",
			Tx:          message,
		}
		if err := c.client.CallContext(c.ctx, &reply, "wallet_cKBSignMsg", param); err != nil {
			return nil, fmt.Errorf("remoteRpcClient.Call err: %s", err.Error())
		}
		if reply.Errno == 0 {
			signTxStr := reply.Data.(string)
			signTxBys, err := hex.DecodeString(signTxStr)
			if err != nil {
				return nil, fmt.Errorf("hex.DecodeString signed tx err: %s", err.Error())
			}
			return signTxBys, nil
		} else {
			return nil, fmt.Errorf("remoteRpcClient.Call err: %s", reply.Errmsg)
		}
	}
	return nil, fmt.Errorf("sign func is nil")
}

type HandleSignCkbMessage func(message string) ([]byte, error)
