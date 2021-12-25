package controller

import (
	"awesomeProject/Util"
	"awesomeProject/common"
	"awesomeProject/dto"
	"awesomeProject/model"
	"awesomeProject/response"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	//获取参数
	db := common.GetDB()
	name := c.PostForm("Name")
	Mail := c.PostForm("Mail")
	password := c.PostForm("Password")
	if isEmailExist(db, Mail) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该邮箱已存在")
		return
	}
	if !Util.VerifyEmailFormat(Mail) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "请输入有效邮箱地址")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码长度应大于6位")
		return
	}
	//没有给名称随机取名
	if len(name) == 0 {
		name = Util.RandomString(10)
	}
	hasePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "加密失败")
		return
	}

	//判断邮箱号是否存在
	newUser := model.User{
		Name:      name,
		Mail:      Mail,
		Password:  string(hasePassword),
		Privilege: "01",
	}
	db.Create(&newUser)

	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	response.Success(c, gin.H{"token": token}, "注册成功")
}

func isEmailExist(db *gorm.DB, Mail string) bool {
	var user model.User
	db.Where("Mail = ?", Mail).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

func Login(c *gin.Context) {
	db := common.GetDB()
	//获取数据
	//使用map获取请求参数
	var requestUser = model.User{}
	c.Bind(&requestUser)

	//获取参数
	Mail := requestUser.Mail
	password := requestUser.Password
	//数据验证
	if !isEmailExist(db, Mail) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该邮箱不存在")
		return
	}
	if !Util.VerifyEmailFormat(Mail) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "请输入有效邮箱地址")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码长度应大于6位")
		return
	}

	//判断手机号是否存在
	var user model.User
	db.Where("Mail = ?", Mail).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}

	//返回结果
	response.Success(c, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}
