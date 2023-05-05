package handler

import (
	"fmt"
	"moratal/apirouter/common"
)

// 存储生成对应的方法
// key method-reqPath string format='Method-reqPath', e.g. key='POST-/api/nodes/v1/activate'
//
//	value FuncHandler
var funcMap = make(map[string]FuncHandler)

type RegisterApi struct {
	Method  string `json:"method"`
	Pattern string `json:"pattern"`
	//todo 接口需要细化实现
	Func FuncHandler
}

var apiRouter *common.Router
var registerApi = []RegisterApi{
	//各种路由及方法eg:
	//{Method: http.MethodGet, Pattern: "/api/mortal/v1/new", Func: nil},
}

// MustRegisterApi 注册路由
func MustRegisterApi() {
	router := common.NewRouter("")
	//api List注册
	var apiList []common.ApiInfo
	for _, api := range registerApi {
		apiList = append(apiList, common.ApiInfo{Method: api.Method, Pattern: api.Pattern})
	}
	router.AddRouter(apiList...)
	//初始化路由handler匹配器
	for _, api := range registerApi {
		funcMap[fmt.Sprintf("%s-%s", api.Method, api.Pattern)] = api.Func
	}
}
