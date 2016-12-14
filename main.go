package main

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/mozillazg/request"
	  "golang.org/x/text/encoding/simplifiedchinese"
    "golang.org/x/text/transform"
)

const HTTP_MAX_TRY_TIMES = 3

var(
	baseURL string
	xuehao string
	password string
)

func init(){
	baseURL = "http://202.116.160.166"
	xuehao = "201623251128"
	password = "1296580449xy"
}

func post(url string, data map[string]string, cookies []*http.Cookie)([]byte, []*http.Cookie, error){
	c := new(http.Client)
	req := request.NewRequest(c)
	if data!=nil {// 转化为gbk
		
	}
	req.Data = data
	var err error
	var resp *request.Response
	for hasTryTimes := 1;  hasTryTimes< HTTP_MAX_TRY_TIMES; hasTryTimes++ {// 重试3次
		resp, err = req.Post(url)
		if err==nil {
			defer resp.Body.Close()
			break
		}
	}
	if err!=nil {
		return nil,nil,err
	}
	respBytes, _ := ioutil.ReadAll(transform.NewReader(resp.Body, simplifiedchinese.GBK.NewEncoder()))
	return respBytes, resp.Cookies(), nil
}

func loginIn()(map[string]string, error){
	var loginUrl string
	loginUrl = baseURL + "/default_ysdx.aspx"
	loginData := map[string]string{
		"TextBox1": xuehao,
		"TextBox2":   password,
		"RadioButtonList1":   "学生",
		"Button1":   "登录",
	}
	 respBytes,cookies, err:= post(loginUrl, loginData, nil)
	 if err!=nil {
		 fmt.Println("error:",err.Error())
	 }
	 fmt.Println("resp:",string(respBytes))
	 fmt.Println("cookies:",cookies)
	return nil,nil
}


func main()  {
	_, err :=loginIn()
	if err!=nil {
		fmt.Println(err.Error())
	}
	// fmt.Println("hello world");
}