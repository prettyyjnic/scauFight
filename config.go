package scauFight

import (
	"flag"

	"github.com/larspensjo/config"
)

const HTTP_MAX_TRY_TIMES = 3

var zhengFang struct {
	baseURL        string
	loginURL       string
	codeURL        string
	chineseURL     string
	publicClassURL string
	mainURL        string
}
var Config *config.Config

// 正方返回错误的 bytes 数组
var unCorrectRequestBytes = [][]byte{
	[]byte("三秒防刷"),
	[]byte("出错啦"),
	[]byte("请重新登陆"),
	[]byte("登录失败"),
	[]byte("Object moved"),
	[]byte("Service Unavailable"),
	[]byte("Location: /logout.aspx"),
}

// 超时
var timeoutRequestStrings = []string{
	"Object moved",
	"Location: /logout.aspx",
	"请重新登陆",
	"登录失败",
}

// 上课时间冲突或者超过限制
var limitRequestStrings = []string{
	"上课时间冲突",
	"选课门数超过限制",
}

func init() {
	var err error
	configFile := flag.String("configfile", "./config.ini", "General configuration file")
	Config, err = config.ReadDefault(*configFile)
	if err != nil {
		panic("配置文件不存在")
	}
	zhengFang.baseURL, err = Config.String("system", "baseURL")
	if err != nil {
		panic("请设置baseURL!")
	}

	zhengFang.loginURL = zhengFang.baseURL + "/default2.aspx"
	zhengFang.mainURL = zhengFang.baseURL + "/xs_main.aspx?xh="
	zhengFang.codeURL = zhengFang.baseURL + "/CheckCode.aspx"
	zhengFang.chineseURL = zhengFang.baseURL + "/xf_xstyxk_qtk.aspx?xh="
	zhengFang.publicClassURL = zhengFang.baseURL + "/xf_xsqxxxk.aspx?xh="
}
