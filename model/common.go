package model

import "time"

type Result struct {
	Success bool        `json:"success" format:"bool"`                // 是否请求成功
	Code    int         `json:"code" format:"int"`                    // 响应编码
	Message string      `json:"message" format:"string" default:"ok"` // 消息
	Time    string      `json:"time"  format:"string"`                // 请求时间
	Data    interface{} `json:"data" format:"object"`                 // 响应结果
}

type HTTPError struct {
	Success bool   `json:"success" format:"bool" default:"false"`  // 是否请求成功
	Code    int    `json:"code" format:"int" default:"-1"`         // 响应编码
	Message string `json:"message" format:"string" default:"请求失败"` // 消息
}

func BuildSuccess(data interface{}) Result {
	result := Result{}
	result.Code = 0
	result.Success = true
	result.Message = "ok"
	format := time.Now().Format("2006-01-02 15:04:05")
	result.Time = format
	result.Data = data
	return result
}

func BuildError(errMsg string) Result {
	result := Result{}
	result.Message = errMsg
	result.Code = -1
	format := time.Now().Format("2006-01-02 15:04:05")
	result.Time = format
	result.Success = false
	return result
}
