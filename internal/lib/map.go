package lib

import (
	"math"
		)

//返回单位为：千米
func GetDistance(lat1, lat2, lng1, lng2 float64) float64 {
	radius := 6371000.0 //6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	rest:=dist * radius / 1000
	//println(math.Trunc(rest*1e2+0.5) * 1e-2)
	//result2:=strconv.ParseFloat(fmt.Sprintf("%.3f",rest), 3)
	pow10_n := math.Pow10(2)
	res2:=math.Trunc((rest+0.5/pow10_n)*pow10_n) / pow10_n
	return res2
	//return strconv.ParseFloat(fmt.Sprintf("%.3f",rest), 3)
}
