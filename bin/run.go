package main

import (
	"log"

	"github.com/prettyyjnic/scau-fight"
)

func main() {
	// str := `<input id=".*" type="checkbox" name="(.*)">`
	// reg := regexp.MustCompile(`<input id=".*" type="checkbox" name="(.*)">`)
	// matches := reg.FindSubmatch([]byte(str))
	// fmt.Println(matches[1])

	xuehao, err := scauFight.Config.String("student", "xuehao")
	if err != nil {
		panic("学号配置错误！")
	}
	password, err := scauFight.Config.String("student", "password")
	if err != nil {
		panic("密码配置错误！")
	}

	student := scauFight.NewStudent(xuehao, password)
	fightClassA(student)
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

func fightWithClassName(student *scauFight.StudentStruct) {
	className := "大学语文"
	teacherName := "杨汤琛"
	response, err := student.FightChineseClassByClassName(className, teacherName, "")
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(string(response))
	}
}

func fightClassA(student *scauFight.StudentStruct) {
	className := "中国哲学智慧与现代企业管理"
	response, err := student.FightPublicClassByClassInfo(className)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(string(response))
	}
}
