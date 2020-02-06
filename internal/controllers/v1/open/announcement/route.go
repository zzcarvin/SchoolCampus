package announcement

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party){
	party.Get("s",getAnnouncement)
}
