package account

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

type roleback struct {
	Id        int    `json:"id" xorm:"id pk"`
	Name      string `json:"name" xorm:"name"`
	Privilege int    `json:"privilege" xorm:"privilege"`
	Avatar    string `json:"avatar" xorm:"avatar"`
	Roles     string `json:"roles" xorm:"roles"`
	Introduce string `json:"introduce" xorm:"introduce"`
	Username  string `json:"username" xorm:"username"`
}
type oldpassword struct {
	Password string `json:"password"      xorm:"password"       validate:"lte=128"`
}

func account(ctx iris.Context) {
	account := models.Account{} //用于存放lastlogin

	roleback := roleback{}

	//token := ctx.GetHeader("Authorization")
	////t	:= strings.Split(token,"Bearer")
	//
	//t :=(token[7:len(token)])
	//
	//
	// id, err := lib.ParseToken(string(t),[]byte(configs.Conf.Jwt.Key))
	// fmt.Println("id=",id)
	//if err != nil {
	//	lib.NewResponseFail(1,"根据token未查出account id")
	//}

	token := ctx.Values().Get("jwt").(*jwt.Token)
	id := token.Claims.(jwt.MapClaims)["id"]
	fmt.Println("id=", id)
	b, err := lib.Engine.Table("account").
		Join("INNER", "role", "role.id = account.role_id").
		Where("account.id=?", id).
		Cols("role.name", "role.privilege", "account.avatar", "role.introduce", "role.roles", "account.Username", "account.id").
		Get(&roleback)
	fmt.Println("token=", ctx.Request())
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户"))
		return
	}

	res, err1 := lib.Engine.Table("account").ID(id).Update(account)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(roleback))
	lib.NewResponseOK(res)

}
func modifypassword(ctx iris.Context) {
	token := ctx.Values().Get("jwt").(*jwt.Token)
	id := token.Claims.(jwt.MapClaims)["id"]

	fmt.Println("\n\n\n\n用户id\n\n", id)
	//id := ctx.Values().Get("jwt").(*jwt.Token).Claims.(jwt.MapClaims)["id"]
	passwd := models.Password{}
	oldpasswd := oldpassword{}
	account := models.Account{}
	//id := ctx.Params().GetUint64Default("id", 0)
	err := ctx.ReadJSON(&passwd)
	if err != nil {
		lib.NewResponseFail(1, err.Error())
		return

	}

	b, err := lib.Engine.Table("account").Where("id=?", id).Cols("password").Get(&oldpasswd)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该用户密码 "))
		return
	}
	//judge :=strings.EqualFold(oldpasswd.Password,passwd.Oldpassword)
	judge := oldpasswd.Password == passwd.Oldpassword
	if judge == false {
		ctx.JSON(lib.NewResponseFail(1, "旧密码不正确"))
		return
	}
	account.Password = passwd.Newpassword

	res, err := lib.Engine.Table("account").ID(id).Update(&account)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	fmt.Println("受影响的行数", res)
	ctx.JSON(lib.NewResponseOK("修改密码成功"))

}
