package scauFight

import (
	"log"
	"net/http"
	"strings"
)

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

func (student *StudentStruct) LoginOut() {
	student.cookies = nil
}

func (student *StudentStruct) FightPublicClassAuto(courses []*CourseInfo) {
	for _, courseInfo := range courses {
		for {
			if _, err := student.FightPublicClassByClassInfo(courseInfo.CourseName, courseInfo.TeacherName, courseInfo.CourseTime); err != nil {
				log.Println("error :", err.Error())
				if strings.Contains(err.Error(), "上课时间冲突") || strings.Contains(err.Error(), "选课门数超过限制") {
					break
				}
				if strings.Contains(err.Error(), "请重新登陆") || strings.Contains(err.Error(), "Object moved") || strings.Contains(err.Error(), "logout.aspx") || strings.Contains(err.Error(), "登录失败") {
					student.LoginOut()
					log.Println("退出登录")
				}
			} else {
				log.Println("抢课" + courseInfo.CourseName + "成功！")
				break
			}
		}
	}
}
