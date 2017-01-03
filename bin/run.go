package main

import (
	"log"

	"github.com/prettyyjnic/scau-fight"
)

func main() {

	xuehao, err := scauFight.Config.String("student", "xuehao")
	if err != nil {
		panic("学号配置错误！")
	}
	password, err := scauFight.Config.String("student", "password")
	if err != nil {
		panic("密码配置错误！")
	}
	// className := "中国哲学智慧与现代企业管理"

	courses := []*scauFight.CourseInfo{
		&scauFight.CourseInfo{
			CourseName: "植物源化学物质及其应用",
			CourseTime: "",
		},
		&scauFight.CourseInfo{
			CourseName: "丝绸文化(A系列)",
			CourseTime: "周三第9,10节",
		},
		&scauFight.CourseInfo{
			CourseName: "生物安全",
			CourseTime: "周二第9,10节",
		},
		&scauFight.CourseInfo{
			CourseName: "花粉的功能与应用",
			CourseTime: "周四第11,12",
		},
		&scauFight.CourseInfo{
			CourseName: "微量元素与健康",
			CourseTime: "",
		},
	}

	student := scauFight.NewStudent(xuehao, password)
	student.FightPublicClassAuto(courses)
}

func fightWithCode(student *scauFight.StudentStruct) {
	// 抢中文课,修改该字段为要抢的课程的code（右键审查元素查看要选的课程的checkbox的name ）
	courseCode := "kcmcGrid:_ctl4:xk"
	response, err := student.FightChineseClassByClassCode(courseCode)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(string(response))
	}
}

func fightWithClassName(student *scauFight.StudentStruct) error {
	className := "大学语文"
	teacherName := "杨汤琛"
	_, err := student.FightChineseClassByClassName(className, teacherName, "")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("抢课成功！")
	}
	return err
}

func fightClassA(student *scauFight.StudentStruct, className string, courseTime string) error {

	_, err := student.FightPublicClassByClassInfo(className, "", courseTime)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("抢课成功！")
	}
	return err
}
