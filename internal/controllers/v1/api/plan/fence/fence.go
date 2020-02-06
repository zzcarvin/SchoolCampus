package fence

import (
	"github.com/kataras/iris"
	"Campus/internal/models"
	"Campus/internal/lib"
	"strings"
	"github.com/go-xorm/builder"
)
// swagger:parameters  FenceCreateRequest
type FenceCreateRequest struct {
	// in: body
	Body models.PlanFence
}

// 响应结构体
//
// swagger:response    FenceCreateResponse
type FenceCreateResponse struct {
	// in: body
	Body	fenceresponseMessage

	}
type fenceresponseMessage struct {
	models.ResponseType
	Data   models.PlanFence
}


func create(ctx iris.Context) {
	// swagger:route POST /api/plan/fence fence FenceCreateRequest
	//
	// 创建围栏
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: FenceCreateResponse
	planFence := models.PlanFence{}
	err := ctx.ReadJSON(&planFence)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))

	}

	//ctx.JSON(lib.NewResponseOK(planFence))

	res, err := lib.Engine.Table("plan_fence").Insert(&planFence)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(planFence))
}

// swagger:route DELETE /api/plan/fence fence FenceDelete
//
//	 删除围栏
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
	planFence := models.PlanFence{}
	affected, err := lib.Engine.Table("plan_fence").ID(id).Delete(&planFence)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}



// swagger:parameters  FenceUpdateRequest
type FenceUpdateRequest struct {
	// in: body
	Body models.PlanFence
}

// 响应结构体
//
// swagger:response    FenceUpdateResponse
type FenceUpdateResponse struct {
	// in: body
	Body	fenceresponseMessage

}


func update(ctx iris.Context) {
	// swagger:route put /api/plan/fence/:id fence FenceUpdateRequest
	// 修改围栏
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//     Responses:
	//       200: FenceUpdateResponse

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planFence := models.PlanFence{}

	//解析student
	err := ctx.ReadJSON(&planFence)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan_fence").ID(id).Update(&planFence)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

// swagger:route GET /api/plan/fence/:id fence FenceGet
//
// 查询围栏
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     Responses:
//       200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planFence := models.PlanFence{}
	//根据id查询
	b, err := lib.Engine.Table("plan_fence").ID(id).Get(&planFence)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该围栏"))
		return
	}
	ctx.JSON(lib.NewResponseOK(planFence))
}
// swagger:route GET /api/plan/fences fence FenceSearch
//
// 查询围栏
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     Responses:
//       200: Response
func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("plan_fence")

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
	var planfence []models.PlanFence
	err := query.Find(&planfence)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(planfence))
}


