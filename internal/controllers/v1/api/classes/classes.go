package classes

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

// swagger:parameters  ClassesCreateRequest
type ClassesCreateRequest struct {
	// in: body
	Body models.Classes
}

// 响应结构体
//
// swagger:response    ClassesCreateResponse
type ClassesCreateResponse struct {
	// in: body
	Body classesresponseMessage
}
type classesresponseMessage struct {
	// Required: true
	models.ResponseType
	Data models.Classes
}

func create(ctx iris.Context) {
	//swagger:route POST /api/classes classes ClassesCreateRequest
	//
	//创建班级
	//    Consumes:
	//    - application/json
	//
	//    Produces:
	//    - application/json
	//
	//    Responses:
	//    - 200: ClassesCreateResponse
	classes := models.Classes{}

	//classes
	err := ctx.ReadJSON(&classes)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}
	//插入数据
	res, err := lib.Engine.Table("classes").Insert(&classes)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(classes))

}

// swagger:route DELETE /api/classes/:id classes ClassesDelete
//
// 删除班级
//    Consumes:
//    - application/json
//
//    Produces:
//    - application/json
//
//    Responses:
//    - 200: Response
func remove(ctx iris.Context) {
	id := ctx.Params().GetUint64Default("id", 0)
	classes := models.Classes{}
	affected, err := lib.Engine.Table("classes").ID(id).Delete(&classes)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))

}

// swagger:parameters  ClassesUpdateRequest
type ClassesUpdateRequest struct {
	// in: body
	Body models.Classes
}

// 响应结构体
//
// swagger:response    ClassesUpdateResponse
type ClassesUpdateResponse struct {
	// in: body
	Body classesresponseMessage
}

func update(ctx iris.Context) {
	// swagger:route PUT /api/classes/:id classes ClassesUpdateRequest
	//
	// 修改班级
	//    Consumes:
	//    - application/json
	//
	//    Produces:
	//    - application/json
	//
	//    Responses:
	//    - 200: ClassesUpdateResponse
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	classes := models.Classes{}

	//解析classes
	err := ctx.ReadJSON(&classes)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("classes").ID(id).Update(classes)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

// swagger:route GET /api/classes/:id classes ClassesGet
//
// 获取班级
//    Consumes:
//    - application/json
//
//    Produces:
//    - application/json
//
//    Responses:
//    - 200: Response
func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	classes := models.Classes{}
	//根据id查询
	b, err := lib.Engine.Table("classes").ID(id).Get(&classes)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该班级"))
		return
	}
	ctx.JSON(lib.NewResponseOK(classes))
}

// swagger:route GET /api/classess classes ClassesSearch
//
// 查询多条班级
//    Consumes:
//    - application/json
//
//    Produces:
//    - application/json
//
//    Responses:
//    - 200: Response
func search(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("classes")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("class_type") {
		query.And(builder.Like{"class_type", ctx.URLParam("class_type")})
	}
	if ctx.URLParamExists("departmentid") {
		query.Where("department_id=?", ctx.URLParamInt64Default("departmentid", 0))
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
	var classes []models.Classes
	err1 := query.Find(&classes)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err1.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(classes))

}
func getcount(ctx iris.Context) {
	var departmentid int64
	departmentid = ctx.URLParamInt64Default("departmentid", 0)
	var classcount int64
	count := new(models.Classes)
	classcount, err := lib.Engine.Where("department_id=?", departmentid).Count(count)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(classcount))
}
