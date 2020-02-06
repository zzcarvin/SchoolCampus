package department

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"github.com/kataras/iris"

	"github.com/go-xorm/builder"
	"strings"
)

// swagger:parameters  DepartmentCreateRequest
type DepartmentCreateRequest struct {
	// in: body
	Body models.Department
}

// 响应结构体
//
// swagger:response    DepartmentCreateResponse
type DeparmentCreateResponse struct {
	// in: body
	Body departmentresponseMessage
}
type departmentresponseMessage struct {
	// Required: true
	models.ResponseType
	Data models.Department
}

func create(ctx iris.Context) {
	// swagger:route POST /api/plan/department department DepartmentCreateRequest
	//
	// 创建院系
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: DepartmentCreateResponse

	department := models.Department{}

	//解析department
	err := ctx.ReadJSON(&department)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}
	//插入数据
	res, err := lib.Engine.Table("department").Insert(&department)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(department))

}

func remove(ctx iris.Context) {
	// swagger:route  GET /api/plan/department department DepartmentGet
	//
	// 删除院系
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: Response
	id := ctx.Params().GetUint64Default("id", 0)
	department := models.Department{}
	affected, err := lib.Engine.Table("department").ID(id).Delete(&department)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))

}

// swagger:parameters  DepartmentUpdateRequest
type DepartmentUpdateRequest struct {
	// in: body
	Body models.Department
}

// 响应结构体
//
// swagger:response    DepartmentUpdateResponse
type DepartmentUpdateResponse struct {
	// in: body
	Body departmentresponseMessage
}

func update(ctx iris.Context) {
	// swagger:route PUT /api/plan/department department DepartmentUpdateRequest
	//
	// 修改院系
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: DepartmentUpdateResponse

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	department := models.Department{}

	//解析department
	err := ctx.ReadJSON(&department)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("department").ID(id).Update(department)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK(department))

}

func get(ctx iris.Context) {
	// swagger:route GET /api/plan/department department DepartmentGet
	//
	// 获取部门
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: Response
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	department := models.Department{}
	//根据id查询
	b, err := lib.Engine.Table("department").ID(id).Get(&department)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该部门"))
		return
	}
	ctx.JSON(lib.NewResponseOK(department))
}

func search(ctx iris.Context) {
	// swagger:route GET /api/plan/department department DepartmentSearch
	//
	// 获取部门（按条件）+s
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Responses:
	//       200: Response

	//创建查询Session
	query := lib.Engine.Table("department")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
	}
	if ctx.URLParamExists("department_type") {
		query.And("department_type=?", ctx.URLParam("department_type"))
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
	var department []models.Department
	err := query.Find(&department)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(department))

}
func getcount(ctx iris.Context) {
	var departmentcount int64
	count := new(models.Department)
	departmentcount, err := lib.Engine.Where("id >?", 0).Count(count)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(departmentcount))
}
