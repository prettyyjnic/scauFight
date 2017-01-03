package scauFight

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

import "github.com/PuerkitoBio/goquery"

import "log"

const (
	// AREA_ZHUXIAOQU 主校区
	AREA_ZHUXIAOQU = "1"
	// AREA_DONGQU 东区
	AREA_DONGQU = "2"
	// AREA_QILIN 启林
	AREA_QILIN = "3"
)

// 课程信息
type CourseInfo struct {
	CourseName  string
	CourseTime  string
	TeacherName string
}

// GetChineseClass 获取语文课信息
func (student *StudentStruct) GetChineseClass() ([]byte, error) {
	if !student.isLogin {
		err := student.LoginIn()
		if err != nil {
			return nil, err
		}
	}
	headers := map[string]string{
		"Referer": zhengFang.mainURL + student.xuehao,
	}
	respBytes, _, err := get(student.requestClient, zhengFang.chineseURL+student.xuehao, nil, nil, headers)
	return respBytes, err
}

// FightChineseClassByClassCode 根据课程code抢课
func (student *StudentStruct) FightChineseClassByClassCode(classCode string) ([]byte, error) {
	respBytes, err := student.GetChineseClass()
	if err != nil {
		return nil, err
	}
	return student.fightChineseClass(respBytes, classCode)
}

// FightChineseClassByClassName 根据课程信息抢课, 第一个参数 课程名称， 第二个参数 教师名称， 第三个参数 上课时间
func (student *StudentStruct) FightChineseClassByClassName(args ...string) ([]byte, error) {
	respBytes, err := student.GetChineseClass()
	if err != nil {
		return nil, err
	}
	var className string
	var teacherName string
	var courseTime string
	for i, info := range args {
		switch i {
		case 0:
			className = info
		case 1:
			teacherName = info
		case 2:
			courseTime = info
		default:
			log.Println("仅支持3个参数")
		}
	}
	// 解析获取课程的 code
	classCode, err := getChineseClassCodeByClassInfo(respBytes, className, teacherName, courseTime)
	if err != nil {
		return nil, err
	}
	return student.fightChineseClass(respBytes, classCode)
}

// 发送抢语文课请求
func (student *StudentStruct) fightChineseClass(resp []byte, classCode string) ([]byte, error) {
	var err error
	__VIEWSTATE, __VIEWSTATEGENERATOR := getViewState(resp)
	params := map[string]string{
		"__VIEWSTATE":          string(__VIEWSTATE),
		"__VIEWSTATEGENERATOR": string(__VIEWSTATEGENERATOR),
		"Button1":              "提交",
		classCode:              "on",
	}
	headers := map[string]string{
		"Referer":                   zhengFang.chineseURL + student.xuehao,
		"Upgrade-Insecure-Requests": "1",
	}
	resp, _, err = post(student.requestClient, zhengFang.chineseURL+student.xuehao, params, nil, headers)
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile(`<script language='javascript'>alert\('(.*)'\);</script>`)
	matches := reg.FindSubmatch(resp)

	if len(matches) > 1 {
		return nil, errors.New("抢课失败" + string(matches[1]))
	}
	return resp, nil
}

// GetPublicClass 获取A系列课程信息
func (student *StudentStruct) GetPublicClass() ([]byte, error) {
	if !student.isLogin {
		err := student.LoginIn()
		if err != nil {
			return nil, err
		}
	}
	headers := map[string]string{
		"Referer": zhengFang.mainURL + student.xuehao,
	}
	respBytes, _, err := get(student.requestClient, zhengFang.publicClassURL+student.xuehao, nil, nil, headers)
	return respBytes, err
}

// FightPublicClassByClassInfo 抢公选课,第一个参数 课程名称，第二个参数 教师名称，第三个参数 上课时间，第四个参数 上课校区, 第五个参数 课程归属
func (student *StudentStruct) FightPublicClassByClassInfo(args ...string) ([]byte, error) {
	var className string
	var teacherName string
	var courseTime string
	var courseBelongTo string
	var courseArea string
	courseArea = AREA_DONGQU
	var __VIEWSTATE []byte
	var __VIEWSTATEGENERATOR []byte
	for i, info := range args {
		switch i {
		case 0:
			className = info
		case 1:
			teacherName = info
		case 2:
			courseTime = info
		case 3:
			courseArea = info
		case 4:
			courseBelongTo = info
		default:
			log.Println("仅支持5个参数")
		}
	}
	respBytes, err := student.GetPublicClass() // 获取__VIEWSTATE 和 __VIEWSTATEGENERATOR
	if err != nil {
		return nil, err
	}
	log.Println("获取公选课成功")
	__VIEWSTATE, __VIEWSTATEGENERATOR = getViewState(respBytes)
	// 查找课程
	headers := map[string]string{
		"Referer": zhengFang.publicClassURL + student.xuehao,
	}
	params := map[string]string{
		"TextBox1":             className,
		"Button2":              "确定",
		"ddl_kcxz":             "",
		"ddl_ywyl":             "",
		"ddl_xqbs":             courseArea,
		"ddl_sksj":             courseTime,
		"ddl_kcgs":             courseBelongTo,
		"__VIEWSTATE":          string(__VIEWSTATE),
		"__VIEWSTATEGENERATOR": string(__VIEWSTATEGENERATOR),
	}
	log.Println("查找公选课" + className)
	respBytes, _, err = post(student.requestClient, zhengFang.publicClassURL+student.xuehao, params, nil, headers)
	if err != nil {
		return nil, err
	}
	code, err := getPublicClassCodeByClassInfo(respBytes, className, teacherName, courseTime)
	if err != nil {
		return nil, err
	}
	__VIEWSTATE, __VIEWSTATEGENERATOR = getViewState(respBytes)
	params2 := map[string]string{
		"__VIEWSTATE":          string(__VIEWSTATE),
		"__VIEWSTATEGENERATOR": string(__VIEWSTATEGENERATOR),
		code:      "on",
		"Button1": "提交",
	}
	log.Println("发送选课请求")
	respBytes, _, err = post(student.requestClient, zhengFang.publicClassURL+student.xuehao, params2, nil, headers)
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile(`<script language='javascript'>alert\('(.*)'\);</script>`)
	matches := reg.FindSubmatch(respBytes)

	if len(matches) > 1 {
		return nil, errors.New("抢课失败" + string(matches[1]))
	}
	return respBytes, nil
}

// 根据课程信息查找语文课程code
func getChineseClassCodeByClassInfo(resp []byte, className string, teacherName string, courseTime string) (string, error) {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}
	var classCode string
	document.Find("table.datelist tr").Each(func(i int, tr *goquery.Selection) {
		isMatch := true
		tr.Find("td").Each(func(k int, td *goquery.Selection) {
			if classCode != "" || !isMatch {
				return
			}
			switch k {
			case 0:
				if className != td.First().Text() {
					isMatch = false
					return
				}
			case 1:
				if teacherName != "" && teacherName != td.First().Text() {
					isMatch = false
					return
				}
			case 2:
				if courseTime != "" && !strings.Contains(td.First().Text(), courseTime) {
					isMatch = false
					return
				}
			case 9:
				node := td.First().Find("input").First()
				classCode = node.AttrOr("name", "")
				if classCode != "" {
					log.Println("获取课程成功：", classCode)
				}
			}
		})
	})
	if classCode == "" {
		return classCode, errors.New("找不到课程！")
	} else {
		return classCode, nil
	}
}

// 根据课程信息查找公选课程code
func getPublicClassCodeByClassInfo(resp []byte, className string, teacherName string, courseTime string) (string, error) {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}
	var classCode string
	document.Find("table#kcmcGrid tr").Each(func(i int, tr *goquery.Selection) {
		if classCode != "" {
			return
		}
		tr.Find("td").Each(func(k int, td *goquery.Selection) {
			switch k {
			case 2:
				if className != td.First().Text() {
					classCode = ""
					return
				}
			case 4:
				if teacherName != "" && teacherName != td.First().Text() {
					classCode = ""
					return
				}
			case 5:
				if courseTime != "" && !strings.Contains(td.Text(), courseTime) {
					classCode = ""
					return
				}
			case 0:
				node := td.First().Find("input").First()
				classCode = node.AttrOr("name", "")
			}
		})
	})
	if classCode == "" {
		return classCode, errors.New("找不到课程！")
	} else {
		log.Println("获取课程成功：", classCode)
		return classCode, nil
	}
}
