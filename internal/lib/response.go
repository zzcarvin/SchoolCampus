package lib

type Response struct {
	Code    int         //0：成功，1-n：失败
	Data    interface{} //返回主要数据的结构体
	Message string      //成功或失败的信息
}

type nilStruct struct {
}

func NewResponse(code int, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

func NewResponseOK(data interface{}) Response {
	return Response{
		Code:    0,
		Data:    data,
		Message: "OK",
	}
}

func NewResponseFail(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

type MiniResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func SuccessResponse(data interface{}, message string) MiniResponse {
	return MiniResponse{
		Code:    200,
		Data:    data,
		Message: message,
	}
}

func FailureResponse(data interface{}, message string) MiniResponse {
	return MiniResponse{
		Code:    400,
		Data:    data,
		Message: message,
	}
}

func NilStruct() nilStruct {
	return nilStruct{}
}
