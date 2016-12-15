package scauFight

import "regexp"
import "errors"
import "github.com/PuerkitoBio/goquery"
import "bytes"
import "log"

func (student *StudentStruct) GetChineseClass() ([]byte, error) {
	if student.cookies == nil || len(student.cookies) == 0 {
		err := student.LoginIn()
		if err != nil {
			return nil, err
		}
	}
	headers := map[string]string{
		"Referer": zhengFang.mainURL + student.xuehao,
	}
	respBytes, _, err := get(zhengFang.chineseURL+student.xuehao, nil, student.cookies, headers)
	return respBytes, err
}

func (student *StudentStruct) FightChineseClass(classCode string) ([]byte, error) {
	resp, err := student.GetChineseClass()
	if err != nil {
		return nil, err
	}
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
	resp, _, err = post(zhengFang.chineseURL+student.xuehao, params, student.cookies, headers)
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

func (student *StudentStruct) FightChineseClassByClassName(className string, teacherName string, courseTime string) ([]byte, error) {
	resp, err := student.GetChineseClass()
	if err != nil {
		return nil, err
	}
	// 解析获取课程的 code
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return nil, err
	}
	var classCode string
	document.Find("table.datelist tr.alt").Each(func(i int, tr *goquery.Selection) {
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
				if courseTime != "" && courseTime != td.First().Text() {
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
		return nil, errors.New("找不到课程！")
	}
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
	resp, _, err = post(zhengFang.chineseURL+student.xuehao, params, student.cookies, headers)
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
