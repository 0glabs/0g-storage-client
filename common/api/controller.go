package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const httpStatusCodeInternalError = 600

func Wrap(controller func(c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := controller(c)
		if err != nil {
			switch e := err.(type) {
			case *BusinessError:
				// custom business error
				c.JSON(http.StatusOK, e)
			case validator.ValidationErrors:
				// binding error
				c.JSON(http.StatusOK, ErrValidation.WithData(e))
			default:
				// internal server error
				c.JSON(httpStatusCodeInternalError, ErrInternal.WithData(e))
			}
		} else if result == nil {
			c.JSON(http.StatusOK, ErrNil)
		} else {
			c.JSON(http.StatusOK, ErrNil.WithData(result))
		}
	}
}
