package line

import (
	"github.com/kataras/iris"
	"Campus/internal/lib"
	"Campus/internal/models"
	"strings"
	"github.com/go-xorm/builder"
	fmt "fmt"
)
// swagger:parameters  LineCreateRequest
type LineCreateRequest struct {
	// in: body
	Body models.PlanLine
}

// 响应结构体
//
// swagger:response    LineCreateResponse
type LineCreateResponse struct {
	// in: body
	Body	lineresponseMessage


}
type lineresponseMessage struct {
	// Required: true
	models.ResponseType
	Data    models.PlanLine
}



func create (ctx iris.Context){
	// swagger:route POST /api/plan/line line LineCreateRequest
	// 创建线
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: LineCreateResponse
	planline := models.PlanLine{}



	err :=ctx.ReadJSON(&planline)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1,err.Error()))
	}


	res,err :=lib.Engine.Table("plan_line").Insert(&planline)


	if err != nil {
		ctx.JSON(lib.NewResponseFail(1,err.Error()))

	}
	fmt.Println("plan11111",planline)

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(planline))

}
// swagger:route DELETE /api/plan/line line LineDelete
//
//	 删除线
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
	planline := models.PlanLine{}

	affected, err := lib.Engine.Table("plan_line").ID(id).Delete(&planline)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

// swagger:parameters  LineUpdateRequest
type LineUpdateRequest struct {
	// in: body
	Body models.PlanLine
}

// 响应结构体
//
// swagger:response    LineUpdateResponse
type LineUpdateResponse struct {
	// in: body
	Body	lineresponseMessage


}



func update(ctx iris.Context) {
	// swagger:route  PUT /api/plan/line/:id line LineUpdateRequest
	// 修改线
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: LineUpdateResponse



	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planline := models.PlanLine{}


	err := ctx.ReadJSON(&planline)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan_line").ID(id).Update(&planline)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(planline))

}

// swagger:route GET /api/plan/line/:id line LineGet
//
// 查询线(id)
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	planline := models.PlanLine{}
	//根据id查询
	b, err := lib.Engine.Table("plan_line").ID(id).Get(&planline)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该线"))
		return
	}
	ctx.JSON(lib.NewResponseOK(planline))
}
// swagger:route GET /api/plan/lines  line LineSearch
//
// 查询线(字段)
//     Produces:
//     - application/json
//
//     Responses:
//       200: Response
func search(ctx iris.Context) {

	//创建查询Session指针
	query := lib.Engine.Table("plan_line")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("fenceid"){
		query.And(builder.Like{"fence_id",ctx.URLParam("fenceid")})

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
	var planline []models.PlanLine
	err := query.Find(&planline)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(planline))
}


