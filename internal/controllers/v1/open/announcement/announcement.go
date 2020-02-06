package announcement

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

func getAnnouncement(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("announcement")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("departmentid") {
		query.Where("department_id=?",ctx.URLParamInt64Default("departmentid",0))
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
	size := ctx.URLParamIntDefault("size", 50)
	query.Limit(size, page*size)

	//查询
	var Announcement []models.Announcement
	err1 :=query.Find(&Announcement)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err1.Error()))
		return
	}


	ctx.JSON(lib.NewResponseOK(Announcement))


}
