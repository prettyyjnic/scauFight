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
	enc mahonia.Encoder
	dec mahonia.Decoder
)

func init() {
	enc = mahonia.NewEncoder("GBK")
	dec = mahonia.NewDecoder("GBK")
}

func doRequest(c *http.Client, url string, method string, data map[string]string, cookies []*http.Cookie, headers map[string]string) ([]byte, []*http.Cookie, error) {
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
	for _, v := range unCorrectRequestBytes {
		if bytes.Contains(respBytes, v) {
			return errors.New("正方发生错误" + string(v))
		}
	}

	return nil
}

func get(c *http.Client, url string, data map[string]string, cookies []*http.Cookie, headers map[string]string) ([]byte, []*http.Cookie, error) {
	return doRequest(c, url, "GET", data, cookies, headers)
}

func post(c *http.Client, url string, data map[string]string, cookies []*http.Cookie, headers map[string]string) ([]byte, []*http.Cookie, error) {
	return doRequest(c, url, "POST", data, cookies, headers)
}

func getCode(c *http.Client, codeURL string, cookies []*http.Cookie) (string, error) {
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
	if !strings.Contains(string(respByte), "GIF") {
		return "", errors.New("获取验证码失败")
	}
	ioutil.WriteFile("code.gif", respByte, 0666)
	fmt.Println("Please input your code: ")
	var code string
	fmt.Scanln(&code)
	fmt.Println("验证码:", code)
	return code, nil
}
