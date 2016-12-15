package scauFight

import "net/http"

type StudentStruct struct {
	xuehao               string
	password             string
	__VIEWSTATE          string
	__VIEWSTATEGENERATOR string
	cookies              []*http.Cookie
}

func NewStudent(xuehao string, password string) *StudentStruct {
	stu := &StudentStruct{}
	stu.xuehao = xuehao
	stu.password = password
	return stu
}
