package scauFight

import "regexp"

import "errors"
import "log"
import "strings"

func getViewState(respByte []byte) ([]byte, []byte) {
	reg := regexp.MustCompile(`<input type="hidden" name="__VIEWSTATE" value="(.*)" />`)
	matches := reg.FindSubmatch(respByte)
	var __VIEWSTATE []byte
	if len(matches) > 1 {
		__VIEWSTATE = matches[1]
	}
	reg2 := regexp.MustCompile(`<input type="hidden" name="__VIEWSTATEGENERATOR" value="(.*)" />`)
	matches2 := reg2.FindSubmatch(respByte)
	var __VIEWSTATEGENERATOR []byte
	if len(matches2) > 1 {
		__VIEWSTATEGENERATOR = matches2[1]
	}
	return __VIEWSTATE, __VIEWSTATEGENERATOR
}

func (student *StudentStruct) LoginIn() error {
	// 获取cookie
	resp, _, err := get(student.requestClient, zhengFang.loginURL, nil, nil, nil)
	if err != nil {
		return err
	}
	__VIEWSTATE, _ := getViewState(resp)

	code, err := getCode(student.requestClient, zhengFang.codeURL, nil)
	if err != nil {
		return err
	}
	loginData := map[string]string{
		"txtUserName":      student.xuehao,
		"TextBox2":         student.password,
		"RadioButtonList1": "学生",
		"txtSecretCode":    code,
		"Button1":          "",
		"lbLanguage":       "",
		"hidPdrs":          "",
		"hidsc":            "",
		"__VIEWSTATE":      string(__VIEWSTATE),
		// "__VIEWSTATEGENERATOR": string(__VIEWSTATEGENERATOR),
	}

	respBytes, _, err := post(student.requestClient, zhengFang.loginURL, loginData, nil, nil)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`<script language='javascript' defer>alert\('(.*)'\);`)
	matches := reg.FindSubmatch(respBytes)
	if len(matches) > 0 {
		return errors.New("登录失败" + string(matches[1]))
	}
	if strings.Contains(string(respBytes), "欢迎使用正方教务管理系统！请登录") {
		return errors.New("登录失败!")
	}
	log.Println("登录成功！")
	student.isLogin = true
	return nil
}
