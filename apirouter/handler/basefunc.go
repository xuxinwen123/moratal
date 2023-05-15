package handler

import (
	"github.com/google/uuid"
	"net/http"
	"net/url"
)

type ApiProto struct {
	// 请求体
	Path      string      `json:"path"`
	Pattern   string      `json:"pattern"`
	Method    string      `json:"method"`
	RemoteIp  string      `json:"remote_ip"`
	ReqHeader http.Header `json:"req_header"`
	ReqBody   []byte      `json:"req_body"`
	ReqArgs   url.Values  `json:"req_args"`
	//resp
	StatusCode int         `json:"status_code"` //响应码
	RespHeader http.Header `json:"resp_header"`
	RespBody   []byte      `json:"resp_body"`
	Time       string      `json:"time"`
	RequestId  uuid.UUID   `json:"request_id"`
}
type FuncHandler func(proto ApiProto, pathParam map[string]string) error
