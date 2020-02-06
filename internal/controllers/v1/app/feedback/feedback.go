package feedback

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

type requestCreate struct {
	StudentId int    `json:"student_id"`
	Content   string `json:"content"`
}

func create(ctx iris.Context) {

	feedback := models.Feedback{}

	//解析
	err := ctx.ReadJSON(&feedback)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))

		return
	}
	//插入数据
	affected, err := lib.Engine.Table("feedback").Insert(&feedback)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if affected == 0 {
		ctx.JSON(lib.NewResponseFail(1, "反馈添加失败"))
		return
	}
	ctx.JSON(lib.NewResponseOK(feedback))

}
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	feedback := models.Feedback{}
	//根据id查询
	b, err := lib.Engine.Table("feedback").ID(id).Get(&feedback)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该反馈"))
		return
	}
	fmt.Println(feedback)
	ctx.JSON(feedback)
}

func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("feedback")

	//字段查询
	if ctx.URLParamExists("record_id") {
		query.And(builder.Like{"record_id", ctx.URLParam("record_id")})
	}

	//排序
	if ctx.URLParamExists("sort") {
		sort := ctx.URLParam("sort")
		order := strings.ToLower(ctx.URLParamDefault("order", "asc"))
		switch order {
		case "asc":
			query.Asc(sort)
			break
		case "desc":
			query.Desc(sort)
			break
		default:
			ctx.JSON(lib.NewResponseFail(1, "order参数错误，必须是asc或desc"))
			return
		}
	}

	//分页
	page := ctx.URLParamIntDefault("page", 0)
	size := ctx.URLParamIntDefault("size", 0)
	query.Limit(size, page*size)

	//查询
	var feedback []models.Feedback
	err := query.Find(&feedback)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(feedback))
}
