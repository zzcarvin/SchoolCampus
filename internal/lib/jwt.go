package lib

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"errors"
	"Campus/configs"
)

func GenerateJwt(id uint64, key string) (tokenString string, err error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)

	//过期时间,
	expire,_:= time.ParseDuration(configs.Conf.Jwt.Expire)
	claims["exp"] = time.Now().Add(expire).Unix()
	//签发时间
	claims["iat"] = time.Now().Unix()
	//用户id
	claims["id"] = id

	token.Claims = claims

	return token.SignedString([]byte(key))
}

func ParseToken(tokenString string, SecretKey []byte) (id uint64, err error) {
	var token *jwt.Token

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
	id = uint64(claim["id"].(float64))

	//exp = uint64(claim["exp"].(float64))

	return id, nil
}
