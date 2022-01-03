package controller

import (
	"ZJU_BS_Back-End/common"
	"ZJU_BS_Back-End/model"
	"ZJU_BS_Back-End/response"
	"ZJU_BS_Back-End/util"
	_ "database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func GetTasks(c *gin.Context) {
	// query all the tasks
	db := common.GetDB()
	var tasks []model.Assignment
	if err := db.Model(&model.Assignment{}).Order("id desc").Find(&tasks).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	var res []gin.H
	for i := 0; i < len(tasks); i++ {
		res = append(res, gin.H{
			"id":        tasks[i].ID,
			"Uploader":  util.GetUsername(db, tasks[i].UploaderID),
			"PublishAt": tasks[i].CreatedAt.Format("2006-01-02 15:04:05"),
			"Filename":  tasks[i].Filename,
		})
	}
	response.Success(c, gin.H{"data": res}, "所有任务")

}

func GetUnsolvedTasks(c *gin.Context) {
	// query all the tasks remained unsolved
	db := common.GetDB()
	var tasks []model.Assignment
	if err := db.Model(&model.Assignment{}).Where("Annotated != true").Find(&tasks).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	var res []gin.H
	for i := 0; i < len(tasks); i++ {
		res = append(res, gin.H{
			"id":        tasks[i].ID,
			"Uploader":  util.GetUsername(db, tasks[i].UploaderID),
			"PublishAt": tasks[i].CreatedAt.Format("2006-01-02 15:04:05"),
			"Filename":  tasks[i].Filename,
		})
	}
	response.Success(c, gin.H{"data": res}, "未标注任务")
}

func GetAnnotations(c *gin.Context) {
	// query all the Annotations
	db := common.GetDB()
	var tasks []model.Annotation
	if err := db.Model(&model.Annotation{}).Order("id desc").Find(&tasks).Error; err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "查询失败")
		return
	}

	var res []gin.H
	for i := 0; i < len(tasks); i++ {
		var assignment model.Assignment
		if err := db.Model(&model.Assignment{}).Where("id = ?", tasks[i].AssignmentID).First(&assignment).Error; err != nil {
			response.Response(c, http.StatusBadRequest, 400, nil, "查询失败")
			return
		}
		if assignment.Filename == "" {
			response.Response(c, http.StatusBadRequest, 400, nil, "不存在该任务")
			continue
		}
		//需要多返回一个FileName以供管理员检查该Annotation对应的图像
		res = append(res, gin.H{
			"id":           tasks[i].ID,
			"Uploader":     util.GetUsername(db, tasks[i].UploaderID),
			"AssignmentID": tasks[i].AssignmentID,
			"Tags":         tasks[i].Tags,
			"Filename":     assignment.Filename,
		})
	}
	response.Success(c, gin.H{"data": res}, "所有标注")
}

func GetUnsolvedAnnotations(c *gin.Context) {
	db := common.GetDB()
	var tasks []model.Annotation
	if err := db.Model(&model.Annotation{}).Where("Reviewed != true").Order("id desc").Find(&tasks).Error; err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "查询失败")
		return
	}

	var res []gin.H
	for i := 0; i < len(tasks); i++ {
		var assignment model.Assignment
		if err := db.Model(&model.Assignment{}).Where("id = ?", tasks[i].AssignmentID).First(&assignment).Error; err != nil {
			response.Response(c, http.StatusBadRequest, 400, nil, "查询失败")
			return
		}
		if assignment.Filename == "" {
			response.Response(c, http.StatusBadRequest, 400, nil, "不存在该任务")
			continue
		}
		//需要多返回一个FileName以供管理员检查该Annotation对应的图像
		res = append(res, gin.H{
			"id":           tasks[i].ID,
			"Uploader":     util.GetUsername(db, tasks[i].UploaderID),
			"AssignmentID": tasks[i].AssignmentID,
			"Tags":         tasks[i].Tags,
			"Filename":     assignment.Filename,
			"PublishAt":    tasks[i].CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.Success(c, gin.H{"data": res}, "所有标注")
}

func PublishTask(c *gin.Context) {
	db := common.GetDB()
	user, exist := c.Get("user")
	if !exist {
		response.Response(c, http.StatusUnauthorized, 401, nil, "user not found")
		return
	}
	uploaderID := user.(model.User).ID
	FileName := c.PostForm("Filename")

	var newTask = model.Assignment{
		UploaderID: uploaderID,
		Filename:   FileName, //filename includes "pics/"
		Annotated:  false,
		Reviewed:   false,
	}
	db.Create(&newTask)
	response.Success(c, nil, "发布成功")
}

func PublishAnnotation(c *gin.Context) {
	db := common.GetDB()

	user, exist := c.Get("user")
	if !exist {
		response.Response(c, http.StatusUnauthorized, 401, nil, "user not found")
		return
	}
	UploaderID := user.(model.User).ID
	AssignmentID := c.PostForm("AssignmentID")
	Tags := c.PostForm("Tags")
	if Tags == "" {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "标注不能为空")
		return
	}
	assignment := model.Assignment{}
	if err := db.Model(&model.Assignment{}).Where("id = " + AssignmentID).Take(&assignment).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该任务不存在")
		return
	}
	assignment.Annotated = true
	db.Save(&assignment)
	var newAnnotation = model.Annotation{
		UploaderID:   UploaderID,
		AssignmentID: assignment.ID,
		Tags:         Tags,
		Reviewed:     false,
	}
	db.Create(&newAnnotation)

	response.Success(c, nil, "标注完成")
}

func PassAnnotation(c *gin.Context) {
	db := common.GetDB()
	user, exist := c.Get("user")
	if !exist {
		response.Response(c, http.StatusUnauthorized, 401, nil, "user not found")
		return
	}
	ReviewUserID := user.(model.User).ID
	AssignmentID := c.PostForm("AssignmentID")
	AnnotationID := c.PostForm("AnnotationID")

	assignment := model.Assignment{}
	annotation := model.Annotation{}
	if err := db.Model(&model.Annotation{}).Where("id = " + AnnotationID).Take(&annotation).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	if err := db.Model(&model.Assignment{}).Where("id = " + AssignmentID).Take(&assignment).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}

	assignment.Reviewed = true
	assignment.Tags = annotation.Tags
	annotation.Reviewed = true
	annotation.ReviewUserID = ReviewUserID
	db.Save(&assignment)
	db.Save(&annotation)
	response.Success(c, nil, "审核通过")
}

func DeleteAnnotation(c *gin.Context) {
	db := common.GetDB()

	AssignmentID := c.PostForm("AssignmentID")
	AnnotationID := c.PostForm("AnnotationID")

	assignment := model.Assignment{}
	if err := db.Model(&model.Assignment{}).Where("id = " + AssignmentID).Take(&assignment).Error; err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "查询失败")
		return
	}
	if err := db.Delete(&model.Annotation{}, "id = "+AnnotationID).Error; err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "驳回失败")
		return
	}

	assignment.Annotated = false
	response.Success(c, nil, "标注驳回完成")
}
