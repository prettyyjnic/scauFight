package scauFight

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"io"

	"bytes"

	"strings"

	"github.com/axgle/mahonia"
	"github.com/mozillazg/request"
)

var (
	c   *http.Client
	enc mahonia.Encoder
	dec mahonia.Decoder
)

func init() {
	c = new(http.Client)
	enc = mahonia.NewEncoder("GBK")
	dec = mahonia.NewDecoder("GBK")
}

func doRequest(url string, method string, data map[string]string, cookies []*http.Cookie, headers map[string]string) ([]byte, []*http.Cookie, error) {
	req := request.NewRequest(c)
	req.Headers = headers
	req.Data = make(map[string]string)
	if data != nil { // 转化为gbk
		for k, v := range data {
			if k == "TextBox1" { // 特殊处理
				req.Data[k] = v
			} else {
				req.Data[enc.ConvertString(k)] = enc.ConvertString(v)
			}
			// req.Data[k] = enc.ConvertString(v)
		}
	}
	if cookies != nil {
		req.Cookies = make(map[string]string)
		for _, cookie := range cookies {
			req.Cookies[cookie.Name] = cookie.Value
		}
	}
	var err error
	var resp *request.Response
	var respReader io.Reader
	for hasTryTimes := 1; hasTryTimes < HTTP_MAX_TRY_TIMES; hasTryTimes++ { // 重试3次
		if method == "GET" {
			resp, err = req.Get(url)
		} else {
			resp, err = req.Post(url)
		}
		if err == nil {
			defer resp.Body.Close()
			break
		}
	}
	if err != nil {
		return nil, nil, err
	}
	respReader = dec.NewReader(resp.Body)
	respBytes, _ := ioutil.ReadAll(respReader)
	err = checkResult(respBytes)
	if err != nil {
		return nil, nil, err
	}
	return respBytes, resp.Cookies(), nil
}

func checkResult(respBytes []byte) error {
	unCorrectBytes := [][]byte{
		[]byte("三秒防刷"),
		[]byte("出错啦"),
		[]byte("请重新登陆"),
		[]byte("Object moved"),
		[]byte("Service Unavailable"),
		[]byte("Location: /logout.aspx"),
	}
	for _, v := range unCorrectBytes {
		if bytes.Contains(respBytes, v) {
			return errors.New("正方发生错误" + string(v))
		}
	}

	return nil
}

func get(url string, data map[string]string, cookies []*http.Cookie, headers map[string]string) ([]byte, []*http.Cookie, error) {
	return doRequest(url, "GET", data, cookies, headers)
}

func post(url string, data map[string]string, cookies []*http.Cookie, headers map[string]string) ([]byte, []*http.Cookie, error) {
	return doRequest(url, "POST", data, cookies, headers)
}

func getCode(codeURL string, cookies []*http.Cookie) (string, error) {
	req := request.NewRequest(c)
	req.Cookies = make(map[string]string)
	for _, cookie := range cookies {
		req.Cookies[cookie.Name] = cookie.Value
	}
	resp, err := req.Get(codeURL)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	respByte, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("respByte", string(respByte))
	if strings.Contains(string(respByte), "Service Unavailable") {
		return "", errors.New("正方发生错误Service Unavailable")
	}
	ioutil.WriteFile("code.gif", respByte, 0666)
	fmt.Println("Please input your code: ")
	var code string
	fmt.Scanln(&code)
	fmt.Println("验证码:", code)
	return code, nil
}
