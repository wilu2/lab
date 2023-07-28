package middleware

import (
	"net/http"

	"gitlab.intsig.net/textin-gateway/pkg/openapi"

	"github.com/gin-gonic/gin"
)

type ValidatorGetter interface {
	Validator() *openapi.Validator
}

func OpenapiValidationMiddleware(getter ValidatorGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		v := getter.Validator()
		if v != nil {
			err := v.ValidateRequest(c, c.Request)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
					"code": 406,
					"msg":  err.Error(),
				})
				return
			}
		}
		c.Next()
	}
}
