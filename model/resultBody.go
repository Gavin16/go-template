package model

import "time"

type Result struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    string      `json:"time"`
	Version string      `json:"version"`
	Data    interface{} `json:"data"`
}

func BuildSuccess(data interface{}) Result {
	result := Result{}
	result.Code = 200
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
	result.Code = 500
	format := time.Now().Format("2006-01-02 15:04:05")
	result.Time = format
	result.Success = false
	return result
}
