package models

type Response struct {
	Code    int         //0：失败，1：成功
	Data    interface{} //返回主要数据的结构体
	Message string      //成功或失败的信息
}


// 注释所用的父类模型
//
// 返回体
// It's also used as one of main axes for reporting.
//
// A user can have friends with whom they can share what they like.
//
// swagger:model
type ResponseType struct {



		// 请求状态，0成功，1失败
		//
		// Required: true
		// Example: 0
		Code int

		// 提示信息
		//
		// Required: true
		// Example: 成功
		Message string

}

// 注释所用的父类模型
//
// 返回体
// It's also used as one of main axes for reporting.
//
// A user can have friends with whom they can share what they like.
//
// swagger:model
type APPResponseType struct {



	// 请求状态，0成功，1失败
	//
	// Required: true
	// Example: 0
	Code int

	// 提示信息
	//
	// Required: true
	// Example: 成功
	Message string

}