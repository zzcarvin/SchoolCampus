package internal

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
	jwt2 "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/core/errors"
	"Campus/internal/models"
)

//jwt 加盐内容
const (
	SecretKey = "this"
)

type JwtUser struct {
	Id  uint64
	Exp uint64
}

/***********************************************************************************************
*函数名 ： parseToken
*函数功能描述 ： 使用传入的token和“盐”返回解析token的claims
*函数参数 ：tokenString需要解析的jwt，  SecretKey 解析需要的“盐”
*函数返回值 ： claims jwt.Claims解析后返回的payload部分 ，err error 验证错误具体原因，或者验证通过返回nil
*作者 ：Sun
*函数创建日期 ：3.25
*函数修改日期 ：
*修改人 ：
*修改原因 ：
*版本 ：0.0.1
*历史版本 ：0.0.1
***********************************************************************************************/
func ParseToken(tokenString string, SecretKey []byte) (jwtuser JwtUser, err error) {
	var token *jwt.Token
	var jwtUserInfo JwtUser
	token, err = jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	//claims = token.Claims
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}
	//将claims写入新的jwtUser
	//id 返回的是float64
	jwtUserInfo.Id = uint64(claim["id"].(float64))
	jwtUserInfo.Exp = uint64(claim["exp"].(float64))

	return jwtUserInfo, nil
}

//先验证token,验证通过解析token，写入Jwtuser，
func ValidateJWT(ctx context.Context) {

	//返回的结构体
	geturlResponse := models.Response{}
	//返回的message
	message := ""

	//验证配置
	jwtConfig := jwt2.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//加密的“盐”
			return []byte(SecretKey), nil
		},
		Debug: true,
		//签名方法
		SigningMethod: jwt.SigningMethodHS256,
		//过期token，自动return
		Expiration: true,
		//jwt验证错误处理
		ErrorHandler: func(i context.Context, s string) {
			message = "token error"
			return
		},
		Extractor: jwt2.FromAuthHeader,
	}
	//开启验证
	jwt2.New(jwtConfig).Serve(ctx)

	//从请求头获取token
	token, err := jwt2.FromAuthHeader(ctx)
	if err != nil {
		println(err)
		//ctx.WriteString("获取token失败")
		message = "获取token失败"
		geturlResponse = SetResponse(0, nil, message)
		ctx.JSON(geturlResponse)
		return
	}
	//解析返回claims
	jwtUserInfo, err := ParseToken(token, []byte(SecretKey))
	if err != nil {
		println(err)
		message = "解析token失败"
		geturlResponse = SetResponse(0, nil, message)
		ctx.JSON(geturlResponse)
		//这里有一个大问题！！！！本应该在上面出现token error时就return的没有返回，反而往下走了！！！
		//ctx.WriteString("parse token error")
		return
	}
	//获取id

	println(jwtUserInfo.Id)
}
