package code

import (
	"bytes"
	"das-account-indexer/prometheus"
	"encoding/json"
	"fmt"
	api_code "github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var log = logger.NewLogger("api_code", logger.LevelDebug)

func DoMonitorLog(method string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()
		statusCode := ctx.Writer.Status()

		var resp ApiResp
		if statusCode == http.StatusOK && blw.body.String() != "" {
			if err := json.Unmarshal(blw.body.Bytes(), &resp); err != nil {
				log.Warn("DoMonitorLog:", method, err.Error())
			}
			if resp.ErrNo != api_code.ApiCodeSuccess {
				log.Warn("DoMonitorLog:", method, resp.ErrNo, resp.ErrMsg)
			}
		}
		if resp.ErrNo == api_code.ApiCodeSuccess {
			resp.ErrMsg = ""
		}
		prometheus.Tools.Metrics.Api().WithLabelValues(method, fmt.Sprint(statusCode), fmt.Sprint(resp.ErrNo), resp.ErrMsg).Observe(time.Since(startTime).Seconds())
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyWriter) Write(bys []byte) (int, error) {
	b.body.Write(bys)
	return b.ResponseWriter.Write(bys)
}
