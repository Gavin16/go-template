package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

/**
 * 使用 resty 库实现http请求发送
 */

func Get(url string, headers map[string]string) (*resty.Response, error) {
	return DoGet(url, headers, nil)
}

func GetWithParam(url string, headers map[string]string) (*resty.Response, error) {
	return DoGet(url, headers, nil)
}

func DoGet(url string, headers map[string]string, param map[string]interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}

	result, err := req.Get(url)
	if err != nil {
		fmt.Println("get request error:", err)
		return nil, err
	}
	return result, nil
}

func Post(url string, param map[string]interface{}) (*resty.Response, error) {
	headers := make(map[string]string)
	headers["content-type"] = "application/json"
	return DoPost(url, headers, param)
}

// PostFile POST文件上传 files传参 => 文件名:[param名:文件内容]
func PostFile(url string, headers map[string]string, files map[string]map[string][]byte) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	if len(files) > 0 {
		for fileName, paramKV := range files {
			for paramName, fileBytes := range paramKV {
				req.SetFileReader(paramName, fileName, bytes.NewReader(fileBytes))
			}
		}
	}
	result, err := req.Post(url)
	if err != nil {
		fmt.Println("file post request error:", err)
		return nil, err
	}
	return result, nil
}

// PostMulti POST请求发送 multipart/form-data 数据
func PostMulti(url string, headers map[string]string, files map[string]map[string][]byte,
	params map[string]string) (*resty.Response, error) {

	client := resty.New()
	req := client.R()
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	if len(files) > 0 {
		for fileName, paramKV := range files {
			for paramName, fileBytes := range paramKV {
				req.SetFileReader(paramName, fileName, bytes.NewReader(fileBytes))
			}
		}
	}
	if len(params) > 0 {
		req.SetFormData(params)
	}

	result, err := req.Post(url)
	if err != nil {
		fmt.Println("file post request error:", err)
		return nil, err
	}
	return result, nil
}

// DoPost 模式使用json 格式作为传参
func DoPost(url string, headers map[string]string, param map[string]interface{}) (*resty.Response, error) {
	client := resty.New()
	req := client.R()
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}

	if len(param) > 0 {
		marshal, err := json.Marshal(param)
		if err != nil {
			fmt.Println("param marshal error:", err)
			return nil, err
		}
		req.SetBody(string(marshal))
	}
	result, err := req.Post(url)
	if err != nil {
		fmt.Println("common post request error:", err)
		return nil, err
	}
	return result, nil
}
