// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://localhost:8000/swagger/index.html",
        "contact": {},
        "license": {
            "name": "Apache 2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/demo/getUserById": {
            "get": {
                "description": "get user info by id",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "demo样例"
                ],
                "summary": "get user by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "500": {
                        "description": "内部错误",
                        "schema": {
                            "$ref": "#/definitions/model.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/demo/hello": {
            "get": {
                "description": "say hello to given name",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "demo样例"
                ],
                "summary": "say hello",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功",
                        "schema": {
                            "$ref": "#/definitions/model.Result"
                        }
                    },
                    "500": {
                        "description": "内部错误",
                        "schema": {
                            "$ref": "#/definitions/model.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "响应编码",
                    "type": "integer",
                    "format": "int",
                    "default": -1
                },
                "message": {
                    "description": "消息",
                    "type": "string",
                    "format": "string",
                    "default": "请求失败"
                },
                "success": {
                    "description": "是否请求成功",
                    "type": "boolean",
                    "format": "bool",
                    "default": false
                }
            }
        },
        "model.Result": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "响应编码",
                    "type": "integer",
                    "format": "int"
                },
                "data": {
                    "description": "响应结果",
                    "format": "object"
                },
                "message": {
                    "description": "消息",
                    "type": "string",
                    "format": "string",
                    "default": "ok"
                },
                "success": {
                    "description": "是否请求成功",
                    "type": "boolean",
                    "format": "bool"
                },
                "time": {
                    "description": "请求时间",
                    "type": "string",
                    "format": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "age": {
                    "description": "年龄",
                    "type": "integer",
                    "default": 27
                },
                "email": {
                    "description": "邮箱",
                    "type": "string",
                    "default": "JoeOK@gmail.com"
                },
                "id": {
                    "description": "ID",
                    "type": "integer",
                    "format": "int64",
                    "default": 1
                },
                "name": {
                    "description": "姓名",
                    "type": "string",
                    "default": "Joe"
                },
                "nation": {
                    "description": "国籍",
                    "type": "string",
                    "default": "USA"
                },
                "title": {
                    "description": "职位",
                    "type": "string",
                    "default": "engineer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "go-template(Replace with your app name)",
	Description:      "请求状态码定义\ncode= 0, 调用成功\ncode=-1, 系统错误\ncode= 1, 提示接口返回的message\ncode=10, 登录token过期\ncode=20, 接口权限错误,没有权限访问该接口",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
