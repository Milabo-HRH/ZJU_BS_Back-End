package controller

import (
	"ZJU_BS_Back-End/model"
	"ZJU_BS_Back-End/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"regexp"
)

func FileUpload(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		response.Response(c, http.StatusUnauthorized, 401, nil, "user not found")
		return
	}
	username := user.(model.User).Name
	file, err := c.FormFile("file")
	filename := username + "_" + c.PostForm("name")
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "upload failed")
		return
	}
	//check filename
	ptn := `^[a-zA-Z0-9_-]{1,99}(.jpg|.png)$`
	reg := regexp.MustCompile(ptn)
	valid := reg.MatchString(filename)
	if !valid {
		response.Response(c, http.StatusBadRequest, 400, nil, "invalid filename")
		return
	}
	//save to local directory
	uid := uuid.New()
	dst := "pics\\" + uid.String() + filename[len(filename)-4:]
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "upload failed")
		return
	}
	response.Success(c, gin.H{
		"url": dst,
	}, "成功")
}

func GetPicture(c *gin.Context) {
	filename := c.Param("filename")
	c.File("pics/" + filename)
}
