package ibeacon

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

func ibeacon(ctx iris.Context) {
	// swagger:route  GET /api/plan/route  route RouteSearch
	//
	// 查询路径 +s
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: Response

	//创建查询Session指针
	query := lib.Engine.Table("plan_route")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
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
	var PlanRoute []models.PlanRoute
	err := query.Find(&PlanRoute)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(PlanRoute))
}
func create(ctx iris.Context) {
	// swagger:route POST /api/plan/points points PointsCreateRequest
	// 创建点
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: PointsCreateResponse
	planpoints := models.PlanPoints{}

	err := ctx.ReadJSON(&planpoints)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	res, err := lib.Engine.Table("plan_points").Insert(&planpoints)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}
	ctx.JSON(lib.NewResponseOK(res))
}
func searchibeacon(ctx iris.Context) {
	query := lib.Engine.Table("plan_points")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
		if ctx.URLParamExists("fence_id") {
			query.And(builder.Like{"fence_id", ctx.URLParam("fence_id")})

		}
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
	var planpoints []models.PlanPoints
	err := query.Where("name!=''").Find(&planpoints)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(planpoints))

}
func updateibeacon(ctx iris.Context) {

	id := ctx.Params().GetUint64Default("id", 0)

	planpoints := models.PlanPoints{}

	//解析student
	err := ctx.ReadJSON(&planpoints)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan_fence").ID(id).Update(&planpoints)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}
