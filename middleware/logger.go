package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"one-api/common"
)

func SetUpLogger(server *gin.Engine) {
	output := gin.DefaultWriter
	if !common.IsLogLevelEnabled(common.LogLevelInfo) {
		output = gin.DefaultErrorWriter
	}
	server.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: output,
		Formatter: func(param gin.LogFormatterParams) string {
			level := common.LogLevelInfo
			switch {
			case param.StatusCode >= http.StatusInternalServerError:
				level = common.LogLevelError
			case param.StatusCode >= http.StatusBadRequest:
				level = common.LogLevelWarn
			}
			if !common.IsLogLevelEnabled(level) {
				return ""
			}
			var requestID string
			if param.Keys != nil {
				requestID = param.Keys[common.RequestIdKey].(string)
			}
			return fmt.Sprintf("[GIN] %s | %s | %3d | %13v | %15s | %7s %s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				requestID,
				param.StatusCode,
				param.Latency,
				param.ClientIP,
				param.Method,
				param.Path,
			)
		},
	}))
}
