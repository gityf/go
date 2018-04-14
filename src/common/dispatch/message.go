package dispatch

import (
	"net/http"
	"reflect"
	logger "github.com/xlog4go"
)

type Message struct {
	ClientLogId string              //客户端日志id
	LogId       int64               //本地日志id
	Sid         string              //区分业务线
	Source      string              //客户端的来源
	Writer      http.ResponseWriter //http响应object
	Request     *http.Request
	RequestInfo map[string]interface{}
}

/*
func (msg *Message) Dispatch(defaultErrNo int, extra interface{}) Responser {
}
*/

// 事件消息
type EventMessage struct {
	*Message

	EventType       uint64    //消息事件类型，确定执行的操作
	EventTime       int64     //消息事件的发生时间，单位毫秒
	MessageBody     string    //消息体
	EventMsgType    uint64    //收集消息类型
	CallFunc func(msg *EventMessage) Responser
}

func (msg *Message) Response(defaultErrNo int, err error) Responser {
	//返回nil 表示成功
	if err == nil {
		//	logger.Warn("TemporaryLog: DispatchMessage(%v)=%v msg=%+v", msg.RequestInfo, err, *msg)
		return doResponse(msg.ClientLogId, 0, "ok", msg.Writer)
	}

	var resp interface{}
	resp = err
	//有明确的返回类型
	switch resp.(type) {
	case Responser:
		//已经发送响应, 只返回
		return resp.(Responser)
	default:
		logger.Warn("TemporaryLog: DispatchMessage(%v)=%v(%v)", msg.RequestInfo, reflect.TypeOf(resp), resp)
	}

	return doResponse(msg.ClientLogId, -1, err.Error(), msg.Writer)
}

//异步结束, 不返回, 需要监控和打日志
func (msg *Message) NoResponse(defaultErrNo int, err error) Responser {
	r := &HttpResponse{
		ErrNo:  0,
		ErrMsg: "ok",
		LogId:  msg.ClientLogId,
	}

	//返回nil 表示成功
	if err == nil {
		//		logger.Warn("TemporaryLog: DispatchMessage(%v)=%v", msg.RequestInfo, err)
		return r
	}

	r.ErrMsg = err.Error()

	var resp interface{}
	resp = err
	//有明确的返回类型
	switch resp.(type) {
	case Responser:
		//已经发送响应, 只返回
		return resp.(Responser)
		return r
	default:
		logger.Warn("TemporaryLog: DispatchMessage(%v)=%v(%v)", msg.RequestInfo, reflect.TypeOf(resp), resp)
	}

	//只是返回错误
	if defaultErrNo > 0 {
		//使用给定的错误码
		r.ErrNo = defaultErrNo
		return r
	}
	r.ErrNo = -1
	return r
}



func dispatchEventMessage(msg *EventMessage) error {
	return msg.CallFunc(msg)
}

//实现MessageDispatcher接口
func (msg *EventMessage) Dispatch(defaultErrNo int, extra interface{}) (resp Responser) {
	resp = msg.CallFunc(msg)
	return
}
func (msg *EventMessage) DispatchDown(defaultErrNo int, extra interface{}) (resp Responser) {
	var err error
	var ddws *DispatchDownWorkspace
	ddws, err = ddwsMgr.GetDispatchDownWorkspace(msg.Sid)
	err = ddws.SendMessage(msg)
	resp = msg.Message.Response(defaultErrNo, err)
	return
}

func (msg *EventMessage) DispatchHandle(defaultErrNo int, extra interface{}) (resp Responser) {
	err := dispatchEventMessage(msg)
	resp = msg.Message.NoResponse(defaultErrNo, err)
	return
}