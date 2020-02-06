package points

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/kataras/iris"
	"time"

	"github.com/go-xorm/builder"
	"strings"
)

// swagger:parameters  PointsCreateRequest
type PointsCreateRequest struct {
	// in: body
	Body models.PlanPoints
}

// 响应结构体
//
// swagger:response    PointsCreateResponse
type PointsCreateResponse struct {
	// in: body
	Body pointsresponseMessage
}
type pointsresponseMessage struct {
	// Required: true
	models.ResponseType
	Data models.PlanPoints
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
	}

	//查询新插入的点是否重复，重复不允许插入
	oldPoints := models.PlanPoints{}
	//因为数据库查询时单精度和双精度不能查询到点位，使用网上的方法：使用数据大于和小于之间
	exi, err := lib.Engine.Table("plan_points").Where("longitude<?", planpoints.Longitude+0.000001).And("longitude>?", planpoints.Longitude-0.000001).
		And("latitude>?", planpoints.Latitude-0.000001).And("latitude<?", planpoints.Latitude+0.000001).
		And("fence_id=?", planpoints.FenceId).Get(&oldPoints)

	if err != nil {
		_, _ = ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	if exi {
		println("重复点！")
		_, _ = ctx.JSON(lib.NewResponseFail(1, "该点已存在，请勿重复插入！"))
		return
	}

	res, err := lib.Engine.Table("plan_points").Insert(&planpoints)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))

	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(planpoints))
}

// swagger:route DELETE /api/plan/points points PointsDelete
//	 删除点
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func remove(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	planpoints := models.PlanPoints{
		Deleted:  1,
		DeleteAt: time.Now(),
	}
	affected, err := lib.Engine.Table("plan_points").ID(id).Update(&planpoints)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

// swagger:parameters    PointsUpdateRequest
type PointsUpdateRequest struct {
	// in: body
	Body models.PlanPoints
}

// 响应结构体
//
// swagger:response    PointsUpdateResponse
type PointsUpdateResponse struct {
	// in: body
	Body pointsresponseMessage
}

func update(ctx iris.Context) {

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
	res, err2 := lib.Engine.Table("plan_points").ID(id).Update(&planpoints)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)
	println("id:", id)
	planpoints := models.PlanPoints{}
	//根据id查询
	b, err := lib.Engine.Table("plan_points").ID(id).Where("deleted=?", 0).Get(&planpoints)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该点"))
		return
	}
	ctx.JSON(lib.NewResponseOK(planpoints))
}

// swagger:route GET /api/plan/points points PointsSearch
// 查询点(字段) +s
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response

func search(ctx iris.Context) {

	//创建查询Session指针
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
	err := query.Where("deleted=?", 0).Find(&planpoints)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(planpoints))

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
