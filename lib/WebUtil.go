package lib

import "github.com/gin-gonic/gin"

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Success(ctx *gin.Context, msg string, data any) {
	ctx.JSON(200, gin.H{"message": msg, "data": data})
}
