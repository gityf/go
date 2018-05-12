package dispatch

import (
	"io"
	"encoding/json"
	"fmt"
	logger "github.com/xlog4go"
	"common/errorcode"
)

type HttpResponse struct {
	ErrNo  int       `json:"errno"`
	ErrMsg string    `json:"errmsg"`
	LogId  string    `json:"logid,omitempty"`
	Data   []*string `json:"data"`
}

type Responser interface {
	//返回错误码, 用于监控
	ErrCode() int
	//返回内容给调用方
	ResponseJson(io.Writer) (int, error)
	//用于打印日志
	String() string
	//继承 error 接口
	Error() string
}

func (r *HttpResponse) ErrCode() int {
	return r.ErrNo
}

func (r *HttpResponse) ResponseJson(w io.Writer) (n int, err error) {
	var s []byte
	var s1 string
	s, err = json.Marshal(r)
	if err != nil {
		logger.Error("json.Marshal err:%v", err)
		s1 = fmt.Sprintf("{\"errno\":%v,\"errmsg\":\"%v\",\"logid\":\"%v\"}", errorcode.ERRNO_JSON_MARSHAL_FAILED, err, r.LogId)
	} else {
		s1 = string(s)
	}
	n, err = io.WriteString(w, s1)
	if err != nil {
		logger.Error("io.WriteString err:%v", err)
	}
	return
}

func (r *HttpResponse) Error() string {
	return fmt.Sprintf("errno=%v,errmsg=%v", r.ErrNo, r.ErrMsg)
}

func (r *HttpResponse) String() string {
	resJson, _ := json.Marshal(r)
	return string(resJson)
}

func doErrorResponse(logid string, errno int, errmsg string, writer io.Writer) Responser {
	return doResponse(logid, errno, errmsg, writer)
}


func doResponse(logid string, errno int, errmsg string, writer io.Writer) (r Responser) {
	r = &HttpResponse{
		ErrNo:  errno,
		ErrMsg: errmsg,
		LogId:  logid,
	}
	_, err := r.ResponseJson(writer)
	if err != nil {
		logger.Error("doResponse err:%v", err)
	}
	return
}