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
	className := "中国哲学智慧与现代企业管理"
	fightClassAUntilSuccess(xuehao, password, className)
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

func fightClassA(student *scauFight.StudentStruct, className string) error {

	_, err := student.FightPublicClassByClassInfo(className)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("抢课成功！")
	}
	return err
}

func fightClassAUntilSuccess(xuehao string, password string, className string) {
	var channel chan int
	var maxChannel = 10
	var successChannel chan int
	successChannel = make(chan int)
	channel = make(chan int, maxChannel)
	for {
		channel <- 1
		go func() {
			defer func() {
				<-channel
			}()
			student := scauFight.NewStudent(xuehao, password)
			err := fightClassA(student, className)
			if err == nil {
				successChannel <- 1
			}
		}()
	}

	for {
		<-successChannel
	}
}
