/**************************************************************
* Overview :
*
* Author   : guozihong
* Created  : 2024-08-21
***************************************************************/

// Package httpclient
package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	jsoniter "github.com/json-iterator/go"
)

type (
	// Param 参数
	Param map[string]interface{}
	// Header 请求头
	Header map[string]string
)

type (
	// PostData post请求数据
	PostData struct {
		BaseUrl  string
		Url      string
		Param    Param
		Header   Header
		UserName string
		Password string
	}
	// GetData get请求数据
	GetData struct {
		BaseUrl  string
		Url      string
		Param    Param
		Header   Header
		UserName string
		Password string
	}
)

// Post 发起POST请求
func Post(data *PostData) (res []byte, code int, err error) {
	// 参数处理
	var buf io.Reader
	if data.Param != nil {
		bytesData, _err := jsoniter.Marshal(data.Param)
		if _err != nil {
			err = _err
			return
		}
		buf = bytes.NewBuffer(bytesData)
	}
	// url拼接
	if data.BaseUrl != "" {
		data.Url = data.BaseUrl + data.Url
	}
	// request创建
	req, err := http.NewRequest("POST", data.Url, buf)
	if err != nil {
		return
	}
	// 设置默认contentType
	req.Header.Set("Content-Type", "application/json")
	return Request(req, data.Header, data.UserName, data.Password)
}

// Get 发起Get请求
func Get(data *GetData) (res []byte, code int, err error) {
	// 参数处理
	if len(data.Param) > 0 {
		param := url.Values{}
		for k, v := range data.Param {
			param.Add(k, fmt.Sprint(v))
		}
		data.Url = data.Url + "?" + param.Encode()
	}
	// url拼接
	if data.BaseUrl != "" {
		data.Url = data.BaseUrl + data.Url
	}
	// request创建
	req, err := http.NewRequest("GET", data.Url, nil)
	if err != nil {
		return
	}

	return Request(req, data.Header, data.UserName, data.Password)
}

// Request 发起POST请求
// resp的body一定要读取完 io.ReadAll()或者 defer io.Copy(io.Discard,resp.Body) 然后Close
// 读取完，不调用Close也没关系，会自动放入连接池
// 不读取完，不调用Close，会导致连接一直保持，会导致goroutine泄露
// 不读取完，调用Close，会导致连接不会放入连接池，每次都会新建连接
func Request(req *http.Request, header Header, userName, password string) (res []byte, code int, err error) {
	// 设置请求头
	for k, v := range header {
		req.Header.Set(k, v)
	}
	// 基础鉴权
	if userName != "" && password != "" {
		req.SetBasicAuth(userName, password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// 关闭
	defer func() {
		_ = resp.Body.Close()
	}()
	res, err = io.ReadAll(resp.Body)
	code = resp.StatusCode

	return
}
