package middleware

import (
	"github.com/kataras/iris"
	jwt2 "github.com/iris-contrib/middleware/jwt"
	"github.com/dgrijalva/jwt-go"
	"Campus/internal/lib"
	"Campus/configs"
)

var JWT *jwt2.Middleware

func JwtInit() {
	cfg := configs.Conf.Jwt

	//验证配置
	jwtConfig := jwt2.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//加密的“盐”
			return []byte(cfg.Key), nil
		},
		Debug: true,
		//解析后的jwt会放置在 ctx.Values().Get("jwt").(*jwt.Token)
		ContextKey: "jwt",
		//签名方法
		SigningMethod: jwt.SigningMethodHS256,
		//过期token，自动return
		Expiration: true,
		//jwt验证错误处理
		ErrorHandler: func(ctx iris.Context, s string) {
			//ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(lib.NewResponseFail(1, s))
			return
		},
		Extractor: jwt2.FromAuthHeader,
	}

	JWT = jwt2.New(jwtConfig)
}
