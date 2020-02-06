package announcement

import (
	"Campus/internal/lib"
	"Campus/internal/models"

	"github.com/go-xorm/builder"

	"github.com/kataras/iris"
	"strings"
	"time"
)

// swagger:parameters  AnnouncementCreateRequest
type AnnouncementCreateRequest struct {
	// in: body
	Body models.Announcement
}

// 响应结构体
//
// swagger:response    AnnouncementCreateResponse
type AnnouncementCreateResponse struct {
	// in: body
	Body AnnouncementresponseMessage
}
type AnnouncementresponseMessage struct {
	// Required: true
	models.ResponseType
	Data models.Announcement
}

func create(ctx iris.Context) {
	//swagger:route POST /api/Announcement Announcement AnnouncementCreateRequest
	//
	//创建公告
	//    Consumes:
	//    - application/json
	//
	//    Produces:
	//    - application/json
	//
	//    Responses:
	//    - 200: AnnouncementCreateResponse
	announcement := models.Announcement{}
	announcement.Apptime = apptime()

	//Announcement
	err := ctx.ReadJSON(&announcement)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return

	}
	//插入数据
	res, err := lib.Engine.Table("announcement").Insert(&announcement)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	lib.NewResponseOK(res)
	ctx.JSON(lib.NewResponseOK( &announcement))

}

// swagger:route DELETE /api/Announcement/:id Announcement AnnouncementDelete
//
// 删除公告
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
	Announcement := models.Announcement{}
	affected, err := lib.Engine.Table("announcement").ID(id).Delete(&Announcement)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))

}

// swagger:parameters  AnnouncementUpdateRequest
type AnnouncementUpdateRequest struct {
	// in: body
	Body models.Announcement
}

// 响应结构体
//
// swagger:response    AnnouncementUpdateResponse
type AnnouncementUpdateResponse struct {
	// in: body
	Body AnnouncementresponseMessage
}

func update(ctx iris.Context) {
	// swagger:route PUT /api/Announcement/:id Announcement AnnouncementUpdateRequest
	//
	// 修改公告
	//    Consumes:
	//    - application/json
	//
	//    Produces:
	//    - application/json
	//
	//    Responses:
	//    - 200: AnnouncementUpdateResponse
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	Announcement := models.Announcement{}

	//解析Announcement
	err := ctx.ReadJSON(&Announcement)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return


	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("announcement").ID(id).Update(Announcement)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))


}

// swagger:route GET /api/Announcement/:id Announcement AnnouncementGet
//
// 获取公告
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
	//
	announcement := models.Announcement{}

	//根据id查询
	b, err := lib.Engine.Table("announcement").ID(id).Get(&announcement)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该公告"))
		return
	}

	ctx.JSON(lib.NewResponseOK(&announcement))
}
func search(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("announcement")

	//字段查询
	if ctx.URLParamExists("name") {
		query.And(builder.Like{"name", ctx.URLParam("name")})
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
	size := ctx.URLParamIntDefault("size", 50)
	query.Limit(size, page*size)

	//查询
	var Announcement []models.Announcement
	err1 := query.Find(&Announcement)
	if err1 != nil {
		ctx.JSON(lib.NewResponseFail(1, err1.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(Announcement))

}

//获取本地时间，并且序列化，便于后期处理
func apptime() string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	//t, _ := time.Parse("2006-01-02", timeStr)
	//fmt.Println(t.Format(time.UnixDate))
	////Unix返回早八点的时间戳，减去8个小时
	//timestamp := t.UTC().Unix() - 8*3600
	//fmt.Println("timestamp:", timestamp)
	return timeStr

}
