package auth

import (
	"github.com/kataras/iris"

	"Campus/configs"
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/satori/go.uuid"
	"time"
	"fmt"
)

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpireIn     int64  `json:"expire_in"`
	//LoginStatus bool `json:"login_status"`
}

type requestLogin struct {
	Username string `json:"username" xorm:"username" validate:"lte=128"`
	Password string `json:"password" xorm:"password" validate:"lte=128"`
}

func login(ctx iris.Context) {
	user := models.Account{}
	user1 := requestLogin{}
	err := ctx.ReadJSON(&user1)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
	}
	b, err := lib.Engine.Table("account").Where("username=?", user1.Username).And("password=?", user1.Password).Get(&user)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户或密码错误"))
		return

	}
	conf := configs.Conf
	tokenString, err := lib.GenerateJwt(uint64(user.Id), conf.Jwt.Key)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, "未获取token"))
		return

	}
	//refreshToken := uuid.Must(uuid.NewV4()).String()
	refreshToken := uuid.NewV4()
	nowTime := time.Now()
	expireIn := nowTime.Add(time.Hour * 24 * 30).Unix()
	token := LoginResponse{
		Token:        tokenString,
		RefreshToken: refreshToken.String(),
		ExpireIn:     expireIn,
	}
	fmt.Printf("token:%v",token)

	ctx.JSON(lib.NewResponseOK(token))

}
