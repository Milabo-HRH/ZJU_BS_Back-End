package middleware

import (
	"ZJU_BS_Back-End/common"
	"ZJU_BS_Back-End/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type form struct {
	Authorization string
}

func AuthMiddleware(ctx *gin.Context) {

	// 获取 authorization header
	tokenString := ctx.GetHeader("Authorization")
	fmt.Println("请求token", tokenString)
	//validate token format
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}
	tokenString = tokenString[7:]
	token, claims, err := common.ParseToken(tokenString)

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}

	//token通过验证, 获取claims中的UserID
	userId := claims.UserId
	DB := common.GetDB()
	var user model.User
	DB.First(&user, userId)

	// 验证用户是否存在
	if user.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}

	//用户存在 将user信息写入上下文
	ctx.Set("user", user)
	ctx.Next()

}

func AuthPrivilege(requiredPrivilege string) gin.HandlerFunc {
	var requiredPrivilegeLevel int
	if requiredPrivilege == "guest" {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}

	if requiredPrivilege == "normal" {
		requiredPrivilegeLevel = 1
	} else if requiredPrivilege == "important" {
		requiredPrivilegeLevel = 2
	} else if requiredPrivilege == "admin" {
		requiredPrivilegeLevel = 3
	} else {
		requiredPrivilegeLevel = 99
	}

	return func(ctx *gin.Context) {
		user, exist := ctx.Get("user")
		if !exist {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "unauthorized",
			})
			ctx.Abort()
			return
		}
		privilege := user.(model.User).Privilege
		privilegeLevel := 0
		if privilege == "normal" {
			privilegeLevel = 1
		} else if privilege == "important" {
			privilegeLevel = 2
		} else if privilege == "admin" {
			privilegeLevel = 3
		}
		if privilegeLevel < requiredPrivilegeLevel {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
