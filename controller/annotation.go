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
			"id":         tasks[i].ID,
			"UploaderID": tasks[i].UploaderID,
			"FileName":   tasks[i].FileName,
			"Tags":       tasks[i].Tags,
			"Reviewed":   tasks[i].Reviewed,
			"Annotated":  tasks[i].Annotated,
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
			"id":       tasks[i].ID,
			"FileName": tasks[i].FileName,
		})
	}
	response.Success(c, gin.H{"data": res}, "未标注任务")

}

func GetAnnotations(c *gin.Context) {
	// query all the Annotations
	db := common.GetDB()
	var tasks []model.Annotation
	if err := db.Model(&model.Annotation{}).Order("id desc").Find(&tasks).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	var res []gin.H
	for i := 0; i < len(tasks); i++ {
		res = append(res, gin.H{
			"id":           tasks[i].ID,
			"UploaderID":   tasks[i].UploaderID,
			"AssignmentID": tasks[i].AssignmentID,
			"Tags":         tasks[i].Tags,
			"Reviewed":     tasks[i].Reviewed,
			"ReviewUserID": tasks[i].ReviewUserID,
		})
	}
	response.Success(c, gin.H{"data": res}, "所有标注")
}

func GetUnsolvedAnnotations(c *gin.Context) {
	db := common.GetDB()
	var tasks []model.Annotation
	if err := db.Model(&model.Annotation{}).Where("Reviewed != true").Order("id desc").Find(&tasks).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}

	var res []gin.H
	for i := 0; i < len(tasks); i++ {
		var FileName string
		if err := db.Model(&model.Assignment{}).Select("FileName").Where("id = ?", tasks[i].AssignmentID).Order("id desc").First(&FileName).Error; err != nil {
			response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
			return
		}
		if FileName == "" {
			response.Response(c, http.StatusUnprocessableEntity, 500, nil, "不存在该任务")
			continue
		}
		//需要多返回一个FileName以供管理员检查该Annotation对应的图像
		res = append(res, gin.H{
			"id":           tasks[i].ID,
			"UploaderID":   tasks[i].UploaderID,
			"AssignmentID": tasks[i].AssignmentID,
			"Tags":         tasks[i].Tags,
			"Reviewed":     tasks[i].Reviewed,
			"ReviewUserID": tasks[i].ReviewUserID,
			"FileName":     FileName,
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
		FileName:   FileName, //filename includes "pics/"
		Annotated:  false,
		Reviewed:   false,
	}
	db.Create(&newTask)
	response.Success(c, nil, "发布成功")
}

func PublishAnnotation(c *gin.Context) {
	db := common.GetDB()
	UploaderID, _ := c.Get("UploaderID")
	AssignmentID, _ := c.Get("AssignmentID")
	Tags, _ := c.Get("Tags")

	assignment := model.Assignment{}
	if err := db.Model(&model.Assignment{}).Where("id = ?", AssignmentID.(int)).Take(&assignment).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "该任务不存在")
		return
	}
	assignment.Annotated = true
	db.Save(assignment)
	var newAnnotation = model.Annotation{
		UploaderID:   UploaderID.(uint),
		AssignmentID: AssignmentID.(int),
		Tags:         Tags.(string),
		Reviewed:     false,
	}
	db.Create(newAnnotation)

	response.Success(c, nil, "标注完成")
}

func PassAnnotation(c *gin.Context) {
	db := common.GetDB()
	AssignmentID, _ := c.Get("AssignmentID")
	AnnotationID, _ := c.Get("AnnotationID")
	ReviewUserID, _ := c.Get("ReviewUserID")
	if !util.VerifyReviewerID(ReviewUserID.(int)) {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "权限不足")
		return
	}
	if !util.VerifyAnnotationID(AnnotationID.(int)) {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "该用户不存在")
		return
	}
	if !util.VerifyAssignmentID(AssignmentID.(int)) {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "该任务不存在")
		return
	}
	//var Tags string
	//if err := db.Model(&model.Annotation{}).Select("Tags").Where("id = ?", AnnotationID).First(&Tags).Error; err != nil {
	//	response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
	//	return
	//}
	//if err := db.Model(&model.Assignment{}).Select("Tags").Where("id = ?", AssignmentID).Update(Tags).Error; err != nil {
	//	response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
	//	return
	//}
	assignment := model.Assignment{}
	annotation := model.Annotation{}
	if err := db.Model(&model.Annotation{}).Where("id = ?", AnnotationID.(int)).Take(&annotation).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	if err := db.Model(&model.Assignment{}).Where("id = ?", AssignmentID.(int)).Take(&assignment).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	assignment.Reviewed = true
	assignment.Tags = annotation.Tags
	annotation.Reviewed = true
	annotation.ReviewUserID = ReviewUserID.(int)
	db.Save(assignment)
	db.Save(annotation)
	response.Success(c, nil, "审核通过完成")
}

func DeleteAnnotation(c *gin.Context) {
	db := common.GetDB()
	AssignmentID, _ := c.Get("AssignmentID")
	AnnotationID, _ := c.Get("AnnotationID")
	ReviewUserID, _ := c.Get("ReviewUserID")
	if !util.VerifyReviewerID(ReviewUserID.(int)) {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "权限不足")
		return
	}
	if !util.VerifyAnnotationID(AnnotationID.(int)) {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "该用户不存在")
		return
	}
	if !util.VerifyAssignmentID(AssignmentID.(int)) {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "该任务不存在")
		return
	}
	assignment := model.Assignment{}
	//annotaton := model.Annotation{}
	if err := db.Delete(&model.Annotation{}, AnnotationID.(int)).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	if err := db.Model(&model.Assignment{}).Where("id = ?", AssignmentID.(int)).Take(&assignment).Error; err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 500, nil, "查询失败")
		return
	}
	assignment.Annotated = false
	response.Success(c, nil, "标注驳回完成")
}
