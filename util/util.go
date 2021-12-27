package util

import (
	"math/rand"
	"regexp"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklqwertyuiopzxcvbnm")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyUploaderID(ID int) bool {
	return true
	//todo: check if the ID is inside the DB and if the ID inside the token
}

func VerifyPictureID(ID int) bool {
	return true
	//todo: check if the picture is inside the DB
}

func VerifyReviewerID(ID int) bool {
	return true
	//todo: check if the reviewer is users with high privilege
}

func VerifyAssignmentID(ID int) bool {
	return true
	//todo: check if the assignmentID is in the DB
}

func VerifyAnnotationID(ID int) bool {
	return true
	//todo: check if the annotationID is in the DB
}
