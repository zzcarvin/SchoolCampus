package auth

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"
	"strconv"
)

type requestExist struct {
	Code string `json:"code" validate:"max=30"`
	Name string `json:"name" validate:"max=15"`
}

type requestBind struct {
	Code      string `json:"code" validate:"max=30"`
	Name      string `json:"name" validate:"max=15"`
	Cellphone string `json:"cellphone" validate:"max=11"`
	Captcha   string `json:"captcha" validate:"len=6"`
}

// 用姓名学号查询学生是否存在，存在返回学生整合信息，不存在直接返回

func exist(ctx iris.Context) {

	//获取学号和姓名
	var faceMaterial requestExist
	if err := ctx.ReadJSON(&faceMaterial); err != nil {
		ctx.JSON(iris.Map{"phone ReadJSON error": "错误的手机号码"})
		//ctx.JSON(lib.FailureResponse())
		fmt.Println("phone ReadJSON error", err)
		return
	}

	//验证
	valPhone := lib.ValidateRequest(faceMaterial)
	if valPhone == false {
		ctx.JSON(lib.NewResponseFail(1, "手机号或状态码格式错误"))
		return
	}

	fmt.Print(faceMaterial.Code, faceMaterial.Name)
	//先查询，存在就将手机号绑定进user，不存在，直接返回
	student := models.Student{}
	b, err := lib.Engine.Table("student").Where("code=?", faceMaterial.Code).And("name=?", faceMaterial.Name).Get(&student)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "学号或姓名不正确，请与你的导员联系"))
		return
	}

	//返回学生的全部信息
	studentAllInfo := models.StudentAllInfos{}
	println("第一次查询学生之后，学生id:", student.Id)
	//根据id查询
	b, err = lib.Engine.Table("student").
		Join("INNER", "classes", "classes.id=student.class_id").
		Join("INNER", "department", "department.id=student.department_id").
		Where("student.id=?", student.Id).
		Cols("student.id", "student.code", "student.name", "classes.name", "department.name", "student.gender", "student.create_at", "student.cellphone").
		Get(&studentAllInfo)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该学生的详细信息"))
		return
	}
	ctx.JSON(lib.NewResponseOK(studentAllInfo))
	return

}

//
// 用姓名学号绑定手机号

func bind(ctx iris.Context) {

	//获取学号和姓名
	faceMaterial := requestBind{}
	if err := ctx.ReadJSON(&faceMaterial); err != nil {
		ctx.JSON(lib.NewResponseFail(1, "手机号或状态码格式错误"))
		fmt.Println("phone ReadJSON error", err)
		return
	}

	//验证
	valPhone := lib.ValidateRequest(faceMaterial)
	if valPhone == false {
		ctx.JSON(lib.NewResponseFail(1, "手机号或状态码格式错误"))
		return
	}

	//兼容接口，如果Captcha不为空，就验证 start
	if faceMaterial.Captcha != "" {
		var redisCaptcha string
		redisConn := lib.GetRedisConn()
		defer redisConn.Close()
		restTimeNew, err := redis.Int(redisConn.Do("TTL", faceMaterial.Cellphone))
		if err != nil {
			ctx.JSON(lib.NewResponseFail(1, "查询验证码过期时间出错"))
			return
		}

		fmt.Println("refresh_token剩余时间：" + strconv.Itoa(restTimeNew))
		//第一步用redis.value来获取数据
		userTimeUnix, err := redis.Values(redisConn.Do("HMGET", faceMaterial.Cellphone, "captcha"))
		if err != nil {
			fmt.Println("redis get failed:", err)
			ctx.JSON(lib.NewResponseFail(1, "从缓存中读取验证码出错"))
			return
		}
		//第二步，用redis.Scan来将数据转换成，我们需要的格式string
		if _, err := redis.Scan(userTimeUnix, &redisCaptcha); err != nil {
			println("scan error ", err)
			ctx.JSON(lib.NewResponseFail(1, "从缓存中读取验证码出错"))
			return
		}
		//把缓存中的验证码与用户的进行对比·
		if redisCaptcha != faceMaterial.Captcha {
			ctx.JSON(lib.NewResponseFail(1, "验证码错误"))
			return
		}
	}
	//验证end

	//插入数据前查询用户手机号是否已绑定学生,后期添加的，防止学生恶意绑定！
	studentExist := models.Student{}
	resExist, err := lib.Engine.Table("student").
		Where("cellphone=?", faceMaterial.Cellphone).
		Exist(&studentExist)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	if resExist == true {
		ctx.JSON(lib.NewResponseOK("该手机号已绑定学生。"))
		return
	}

	//没有绑定的学生才允许绑定手机号
	student := models.Student{}
	student.Cellphone = faceMaterial.Cellphone
	res, err := lib.Engine.Table("student").Where("code=?", faceMaterial.Code).And("name=?", faceMaterial.Name).Update(&student)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	if res == 1 {
		ctx.JSON(lib.NewResponseOK("学生手机号绑定成功"))
		return
	}
	ctx.JSON(lib.NewResponseFail(0, "学生手机号已绑定"))

	println("res lines:", res)
	return
}
