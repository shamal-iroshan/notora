package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func toInt64Strict(ctx *gin.Context, param string) (int64, bool) {
	idStr := ctx.Param(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}
