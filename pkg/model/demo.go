package model

type User struct {
	Id     int64  `json:"id" format:"int64" default:"1"`   // ID
	Name   string `json:"name" default:"Joe"`              // 姓名
	Age    int    `json:"age" default:"27"`                // 年龄
	Nation string `json:"nation" default:"USA"`            // 国籍
	Title  string `json:"title" default:"engineer"`        // 职位
	Email  string `json:"email" default:"JoeOK@gmail.com"` // 邮箱
}
