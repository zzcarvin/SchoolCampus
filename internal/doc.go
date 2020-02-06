// Package classification 学校后台 API.
//
// 学校后台文档
//
//
// 包括验证码，获取学校列表。
//
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host:localhost
//     BasePath:
//     Version: 0.0.1
//     License:
//     Contact: zhouzicheng_@hotmail.com
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package internal

import (


)
//response返回结构体
//
// swagger:response Response
type Response struct {
	// 正常返回
	// in: body
	Body struct {
		// 请求状态，0成功，1失败
		//
		// Required: true
		// Example: 1
		Code int
		// 返回数据
		//
		// Required: true
		// Example:
		Data interface{}
		// 提示信息
		//
		// Required: true
		// Example: 成功
		Message string
	}
}



































































