package run

import (
	"testing"
)

func TestDynamicRoute(t *testing.T) {

	//points, err = DynamicRoute();
	////创建一个请求
	//req, err := http.NewRequest("GET", "/run/health-check", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//// 我们创建一个 ResponseRecorder (which satisfies http.ResponseWriter)来记录响应
	//rr := httptest.NewRecorder()
	//
	////直接使用HealthCheckHandler，传入参数rr,req
	//HealthCheckHandler(rr, req)
	//
	//// 检测返回的状态码
	//if status := rr.Code; status != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v",
	//		status, http.StatusOK)
	//}
	//
	//// 检测返回的数据
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		rr.Body.String(), expected)
	//}
}
