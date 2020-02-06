package ibeacon

import "github.com/kataras/iris"

func RegisterRoutes(party iris.Party){
	party.Get("/getibeacon",ibeacon)
	party.Post("/creatibeacon",create)
	party.Put("/{id:uint64}", updateibeacon)

	//查询带有蓝牙设备名字的地图点位
	party.Get("/getibeaconname",searchibeacon)
}
