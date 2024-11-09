package handlers

import "github.com/gin-gonic/gin"

// JSONE serializes the given error-struct as JSON into the response body.
//
// It also sets the Content-Type as "application/json".
//
// It also attaches the error to the current context, for future handling or logging.
func JSONE(gctx *gin.Context, err error, code int, obj any) {
	gctx.JSON(code, obj)
	_ = gctx.Error(err)
}
