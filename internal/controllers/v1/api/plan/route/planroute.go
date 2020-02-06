package route

import (
	"github.com/kataras/iris"
	"Campus/internal/lib"
	"strings"
	"Campus/internal/models"
	"github.com/go-xorm/builder"
)
// swagger:parameters  RouteCreateRequest
type RouteCreateRequest struct {
	// in: body
	Body models.PlanRoute
}

// 响应结构体
//
// swagger:response    RouteCreateResponse
type RouteCreateResponse struct {
	// in: body
	Body	routeresponseMessage


}
type routeresponseMessage struct {
	// Required: true
	models.ResponseType
	Data    models.PlanRoute
}




func create(ctx iris.Context) {
	// swagger:route POST /api/plan/route route RouteCreateRequest
	//
	// 创建围栏路径
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: RouteCreateResponse


	PlanRoute := models.PlanRoute{}
	err := ctx.ReadJSON(&PlanRoute)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))

	}



	res, err := lib.Engine.Table("plan_route").Insert(&PlanRoute)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(PlanRoute))
}

// swagger:route DELETE /api/plan/route route RouteDelete
//
//	 删除路径
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
	PlanRoute := models.PlanRoute{}
	affected, err := lib.Engine.Table("plan_route").ID(id).Delete(&PlanRoute)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}
// swagger:parameters  RouteUpdateRequest
type RouteUpdateRequest struct {
	// in: body
	Body models.PlanPoints
}

// 响应结构体
//
// swagger:response    RouteUpdateResponse
type RouteUpdateResponse struct {
	// in: body
	Body	routeresponseMessage


}


func update(ctx iris.Context) {
	// swagger:route PUT /api/plan/route/:id  route RouteUpdateRequest
	// 修改路径
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: RouteUpdateResponse
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	PlanRoute := models.PlanRoute{}

	//解析student
	err := ctx.ReadJSON(&PlanRoute)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan_route").ID(id).Update(&PlanRoute)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

// swagger:route GET /api/plan/routes/:id    route RouteGet
//
// 查询路径
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planroute := models.PlanRoute{}
	//根据id查询
	b, err := lib.Engine.Table("plan_route").ID(id).Get(&planroute)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该路径"))
		return
	}
	ctx.JSON(lib.NewResponseOK(planroute))
}
func search(ctx iris.Context) {
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
	size := ctx.URLParamIntDefault("size", 50)
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



