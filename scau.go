package scauFight

import "regexp"

import "errors"

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
	resp, cookies, err := get(zhengFang.loginURL, nil, nil, nil)
	if err != nil {
		return err
	}
	__VIEWSTATE, _ := getViewState(resp)
	student.cookies = cookies
	code, err := getCode(zhengFang.codeURL, student.cookies)
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
	respBytes, _, err := post(zhengFang.loginURL, loginData, student.cookies, nil)
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`<script language='javascript' defer>alert\('(.*)'\);`)
	matches := reg.FindSubmatch(respBytes)
	if len(matches) > 0 {
		return errors.New("登录失败" + string(matches[0]))
	}
	return nil
}
