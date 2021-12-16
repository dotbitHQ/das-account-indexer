package handle

import (
	"context"
	"das-account-indexer/dao"
	"fmt"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/txbuilder"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/scorpiotzh/mylog"
)

var (
	log = mylog.NewLogger("handle", mylog.LevelDebug)
)

type HttpHandle struct {
	Ctx           context.Context
	Red           *redis.Client
	DbDao         *dao.DbDao
	DasCore       *core.DasCore
	TxBuilderBase *txbuilder.DasTxBuilderBase
}

func GetClientIp(ctx *gin.Context) string {
	clientIP := fmt.Sprintf("%v", ctx.Request.Header.Get("X-Real-IP"))
	return fmt.Sprintf("(%s)(%s)", clientIP, ctx.Request.RemoteAddr)
}
