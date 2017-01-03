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
	requestClient        *http.Client
	isLogin              bool
}

func NewStudent(xuehao string, password string) *StudentStruct {
	stu := &StudentStruct{}
	stu.xuehao = xuehao
	stu.password = password
	stu.requestClient = new(http.Client)
	return stu
}

func (student *StudentStruct) LoginOut() {
	student.isLogin = false
}

func (student *StudentStruct) FightPublicClassAuto(courses []*CourseInfo) {
	for _, courseInfo := range courses {
		for {
			if _, err := student.FightPublicClassByClassInfo(courseInfo.CourseName, courseInfo.TeacherName, courseInfo.CourseTime); err != nil {
				log.Println("error :", err.Error())
				var shouldBreak = false
				for i := 0; i < len(limitRequestStrings); i++ {
					if strings.Contains(err.Error(), limitRequestStrings[i]) {
						shouldBreak = true
						break
					}
				}
				if shouldBreak {
					break
				}
				for i := 0; i < len(timeoutRequestStrings); i++ {
					if strings.Contains(err.Error(), timeoutRequestStrings[i]) {
						student.LoginOut()
						log.Println("登录超时")
					}
				}
			} else {
				log.Println("抢课" + courseInfo.CourseName + "成功！")
				break
			}
		}

	}
}
