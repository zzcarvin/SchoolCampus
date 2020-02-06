package internal

import "Campus/internal/models"

//设置response
func SetResponse(code int, dataStruct interface{}, message string) (resData models.Response) {

	Response := models.Response{}

	Response.Code = code
	Response.Data = dataStruct
	Response.Message = message

	return Response

}
