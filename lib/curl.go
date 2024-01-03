package lib

import (
	"comment-chuli/setup"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func CurlReq(method string, url string, timeout time.Duration, body io.Reader, header map[string]string, desc string) (*simplejson.Json, int) {

	client := &http.Client{
		Timeout: timeout, //设置超时
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		setup.Logger.Println(fmt.Sprintf("err,%s创建请求失败:%s", desc, err))
		return nil, http.StatusInternalServerError
	}

	// 设置头信息
	if len(header) > 0 {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		// 检查错误类型，如果是 Timeout，说明发生了超时
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			setup.Logger.Println(fmt.Sprintf("err,%s请求超时:%s", desc, err))
			return nil, http.StatusGatewayTimeout
		} else {
			setup.Logger.Println(fmt.Sprintf("err,%s发生其他错误:%s", desc, err))
			return nil, http.StatusInternalServerError
		}

	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			setup.Logger.Println(fmt.Sprintf("err,%s关闭失败:%s", desc, err))
		}
	}(resp.Body)

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		setup.Logger.Println(fmt.Sprintf("err,%s获取content失败:%s", desc, err))
		return nil, http.StatusInternalServerError
	}

	result, _ := simplejson.NewJson(contents)
	return result, http.StatusOK
}
