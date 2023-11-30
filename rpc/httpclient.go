package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

/**
 * 使用 resty 库实现http请求发送
 * 主要封装 GET,POST,PUT 三种请求方式 application/json, multipart/form-data 两种请求数据类型
 */

type HttpRequest interface {
	Get(param map[string]string) (*resty.Response, error)
	Post(param map[string]string) (*resty.Response, error)
	PostJson(param map[string]interface{}) (*resty.Response, error)
	PostFile(params map[string]string, file map[string][]byte) (*resty.Response, error)
}

// ClientEntity package url+headers in struct
type ClientEntity struct {
	url     string
	headers map[string]string
}

// Get with form-data content type
func (ce ClientEntity) Get(param map[string]string) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(ce.headers) > 0 {
		req.SetHeaders(ce.headers)
	}

	if len(param) > 0 {
		req.SetFormData(param)
	}

	result, err := req.Get(ce.url)
	if err != nil {
		fmt.Println("get request error:", err)
		return nil, err
	}
	return result, nil
}

// Post with form-data content type
func (ce ClientEntity) Post(param map[string]string) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(ce.headers) > 0 {
		req.SetHeaders(ce.headers)
	}

	if len(param) > 0 {
		req.SetFormData(param)
	}

	result, err := req.Post(ce.url)
	if err != nil {
		fmt.Println("post request error:", err)
		return nil, err
	}
	return result, nil
}

// PostJson with application/json content type
func (ce ClientEntity) PostJson(param map[string]interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(ce.headers) > 0 {
		req.SetHeaders(ce.headers)
	}

	if len(param) > 0 {
		marshal, err := json.Marshal(param)
		if err != nil {
			fmt.Println("param marshal error:", err)
			return nil, err
		}
		req.SetBody(string(marshal))
	}

	result, err := req.Post(ce.url)
	if err != nil {
		fmt.Println("postJson request error:", err)
		return nil, err
	}
	return result, nil
}

// PostFile with form-data and file params
func (ce ClientEntity) PostFile(params map[string]string, file map[string][]byte) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(ce.headers) > 0 {
		req.SetHeaders(ce.headers)
	}
	if len(file) > 0 {
		for paramName, fileBytes := range file {
			req.SetFileReader(paramName, paramName, bytes.NewReader(fileBytes))
		}
	}

	if len(params) > 0 {
		req.SetFormData(params)
	}

	result, err := req.Post(ce.url)
	if err != nil {
		fmt.Println("file post request error:", err)
		return nil, err
	}
	return result, nil
}
