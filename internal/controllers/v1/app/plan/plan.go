package plan

import (
	"Campus/internal/lib"
	"Campus/internal/models"
	"fmt"
	"github.com/go-xorm/builder"
	"github.com/kataras/iris"
	"strings"
)

func create(ctx iris.Context) {

	plan := models.Plan{}

	//解析plan
	err := ctx.ReadJSON(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan").Insert(plan)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

func update(ctx iris.Context) {

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	plan := models.Plan{}

	//解析plan
	err := ctx.ReadJSON(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	//TODO 验证数据有效性

	//插入数据
	res, err2 := lib.Engine.Table("plan").ID(id).Update(plan)
	if err2 != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(res))
}

func remove(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	plan := models.Plan{}

	//根据id查询
	affected, err := lib.Engine.Table("plan").ID(id).Delete(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	ctx.JSON(lib.NewResponseOK(affected))
}

func get(ctx iris.Context) {
	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	plan := models.Plan{}
	//根据id查询
	b, err := lib.Engine.Table("plan").ID(id).Get(&plan)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该计划"))
		return
	}
	ctx.JSON(lib.NewResponseOK(plan))
}

func search(ctx iris.Context) {

	//创建查询Session
	query := lib.Engine.Table("plan")

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
	var plan []models.Plan
	err := query.Find(&plan)
	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}

	ctx.JSON(lib.NewResponseOK(plan))
}

//计划规则详情
func planRuleDetail(ctx iris.Context) {

	//取URL参数 id
	id := ctx.Params().GetUint64Default("id", 0)

	plan := models.Plan{}
	//根据id查询
	b, err := lib.Engine.Table("plan").ID(id).Get(&plan)

	if err != nil {
		ctx.JSON(lib.NewResponseFail(1, err.Error()))
		return
	}
	if b == false {
		ctx.JSON(lib.NewResponseFail(1, "未找到该计划"))
		return
	}

	planDate := ""
	planDate = plan.DateBegin.Format("2006年01月02日") + "-" + plan.DateEnd.Format("2006年01月02日")
	println(planDate)
	ctx.Header("content-Type", "text/html; charset=utf-8")
	//ctx.String()
	//minPaceStr := lib.FomPace(float64(plan.MinPace))
	//maxPaceStr := lib.FomPace(float64(plan.MaxPace))
	sumDistance := fmt.Sprintf("%.0f", float32(plan.TotalDistance/1000))
	sumTimes := fmt.Sprintf("%d", plan.TotalTimes)
	//单次最长时间
	maxTimeLong := ""
	if plan.MaxTimeLong != 0 {
		maxTimeLong = `<li><strong>跑步时长</strong>
						<p>单次跑步里程>=` + fmt.Sprintf("%.2f", float32(plan.MinSingleDistance)/1000) + `公里` + "," + fmt.Sprintf("，必须在%.2f小时内完成（包括暂停时间）", float32(plan.MaxTimeLong)/60) + `。超出时间范围不计为有效数据。</p>
					</li>`
	}

	maxWeekTimes := ""
	if plan.MaxWeekTimes != 0 {
		maxWeekTimes = "," + fmt.Sprintf("，单周计成绩总次数为%d次", plan.MaxWeekTimes)
	}

	planRuleDet := `
		<!DOCTYPE html>
		<html>
			<head>
				<meta http-equiv=Content-Type content="text/html;charset=utf-8">
				<title>学期计划规则</title>
				<style type="text/css">
					body {
						font: 400 14px/1.5 'Open Sans', 'PingFang SC', '\5FAE\8F6F\96C5\9ED1', sans-serif;
					}
			
					ul,
					li,
					p {
						margin: 0;
						padding: 0;
					}
			
					ul {
						list-style: none;
					}
			
					ul li {
						line-height: 25px;
						margin-bottom: 20px;
					}
			
					ul li strong {
						font-size: 15px;
						color:#333;
						font-weight: 600;
						display: inline-block;
						margin-bottom: 5px;
					}
			
					p {
						margin-left: 28px;
						color: #666;
					}
			
					li::before {
						content: ".";
						color: #00CAF3;
						font-size: 60px;
						margin-right: 10px;
						display: inline-block;
						vertical-align: middle;
						margin-top: -40px;
					}
				</style>
			</head>
		<body>
			<ul>
				<li><strong>计划日期</strong>
				` + "<p>" + planDate + "</p>" + `
				</li>
				<li><strong>学期目标</strong>
					<p>男生总次数为` + sumTimes + `次，男生总公里数为` + sumDistance + `公里；</p>
					<p>女生总次数为` + sumTimes + `次，女生总公里数为` + sumDistance + `公里；</p>
				</li>
				<li><strong>跑步时段及要求</strong>
					<p>男生规则：里程在0.00~10.00公里，配速在3'0'--9'0',需要经过所有打卡点；</p>
					<p>女生规则：里程在0.00~10.00公里，配速在3'0'--9'0'需要经过所有打卡点；</p>
				</li>
				<li><strong>有效次数上限</strong>
					<p>单日记成绩次数上限为` + fmt.Sprintf("%d次", plan.MaxDayTimes) + maxWeekTimes + `，超过上限有效记录，不会计入总次数。</p>
				</li>
				<li><strong>跑步打卡</strong>
					<p>根据学生的位置随机分配打卡点，到达打卡点时系统会自动打卡；打卡目标完成后，继续跑步，结束跑步时请确保达到合格标准。</p>
				</li>
				` + maxTimeLong + `
				<li><strong>上传时间</strong>
					<p>跑步记录，需在24小时内上传，为有效成绩。</p>
				</li>
			</ul>
		</body>
</html>

`
	ctx.JSON(lib.NewResponse(0, "计划规则详情", planRuleDet))
}

//计划规则详情
func Notify(ctx iris.Context) {
	ctx.Header("content-Type", "text/html; charset=utf-8")

	planRuleDet := `
		<!DOCTYPE html>
		<html>
			<head>
				<meta http-equiv=Content-Type content="text/html;charset=utf-8">
				<title>常见问题解答</title>
				<style type="text/css">
					body {
						font: 400 14px/1.5 'Open Sans', 'PingFang SC', '\5FAE\8F6F\96C5\9ED1', sans-serif;
					}
			
					ul,
					li,
					p {
						margin: 0;
						padding: 0;
					}
			
					ul {
						list-style: none;
					}
			
					ul li {
						line-height: 25px;
						margin-bottom: 20px;
					}
			
					ul li strong {
						font-size: 15px;
						color:#333;
						font-weight: 600;
						display: inline-block;
						margin-bottom: 5px;
					}
			
					p {
						margin-left: 28px;
						color: #666;
					}
			
					li::before {
						content: ".";
						color: #00CAF3;
						font-size: 60px;
						margin-right: 10px;
						display: inline-block;
						vertical-align: middle;
						margin-top: -40px;
					}
				</style>
			</head>
		<body>
			<ul>
				<li><strong>部分IPHONE手机步频为0的问题</strong>
					<p>按以下"设置>应用管理>app>耗电详情>应用启动管理>允许后台活动"路径设置，如果开着有问题，可以尝试关了重开</p>
				</li>
				<li><strong>息屏情况如何跑步</strong>
					<p>打开后台保护系统设置，设置路径参考"应用启动管理>关闭自动管理并打开各允许启动开关"</p>
				</li>
			</ul>
		</body>
</html>

`
	ctx.JSON(lib.NewResponse(0, "常见问题详情", planRuleDet))
}
