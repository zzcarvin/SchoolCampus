package open

import (
	"Campus/internal/controllers/v1/open/announcement"
	"Campus/internal/controllers/v1/open/ibeacon"
	"Campus/internal/controllers/v1/open/mara"
	"github.com/kataras/iris"
	//"Campus/internal/middleware"
	"Campus/internal/controllers/v1/open/jpush"
	//"Campus/internal/controllers/open/messagebus"
	"Campus/internal/controllers/v1/open/share"
)

func RegisterRoutes(app iris.Party) {
	jpush.RegisterRoutes(app.Party("/push"))
	//消息总线
	//messagebus.RegisterRoutes(app.Party("/messagebus"))

	share.RegisterRoutes(app.Party("/share"))
	announcement.RegisterRoutes(app.Party("/announcement"))
	ibeacon.RegisterRoutes(app.Party("/ibeacon"))
	mara.RegisterRoutes(app.Party("/mara"))
}
