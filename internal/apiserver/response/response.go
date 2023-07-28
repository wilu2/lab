package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.intsig.net/textin-gateway/internal/apiserver/code"
	"gitlab.intsig.net/textin-gateway/pkg/validation"
)

// UnifiedResponse 统一返回
type UnifiedResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"msg"`
}

// BadResponse 错误返回
type BadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// HandlerParamsResponse 处理参数响应错误
func HandlerParamsResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, BadResponse{
		Code:    http.StatusBadRequest,
		Message: validation.Error(err),
	})
}

// HandleResponse 统一返回处理 {"msg":"","data":{},"code":200}
func HandleResponse(c *gin.Context, data any, err error) {
	if c.Writer.Status() != 200 { // 手动指定 http Code 如 302 重定向
		return
	}
	if err == nil {
		c.JSON(http.StatusOK, UnifiedResponse{
			Code:    http.StatusOK,
			Data:    data,
			Message: "success",
		})
	} else {
		coder := code.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), BadResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		})
	}
}

func HandleStringResponse(c *gin.Context, data any, err error) {
	if c.Writer.Status() != 200 { // 手动指定 http Code 如 302 重定向
		return
	}
	c.String(http.StatusOK, data.(string))
}

func HandleBinaryResponse(c *gin.Context, fName, fType string, data []byte, err error) {
	if c.Writer.Status() != 200 { // 手动指定 http Code 如 302 重定向
		return
	}

	var fileName string
	fileType := strings.Split(fType, "/")
	if len(fileType) != 2 { // 无文件类型的情况
		fType = "other"
		fileName = fName
	} else {
		if fileType[0] != "image" && fileType[1] != "pdf" { // 非图片和 pdf 的都是 other
			fType = "other"
		}
		fileName = fmt.Sprintf("%s.%s", fName, fileType[1])
	}
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("X-File-Type", fType)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(200, "application/octet-stream", data)
}
