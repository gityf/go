package dispatch

import "time"

/*
func (msg *Message) Dispatch(defaultErrNo int, extra interface{}) Responser {
}
*/

// 事件消息
type EventMessage struct {
	Sid             string
	EventTime       time.Time     //消息事件的发生时间，单位毫秒
	VenderMessage   interface{} // 自定义消息信息
	CallFunc func(msg *EventMessage) Responser
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
	var ddws *DispatchDownWorkspace
	ddws, _ = ddwsMgr.GetDispatchDownWorkspace(msg.Sid)
	_ = ddws.SendMessage(msg)
	return
}

func (msg *EventMessage) DispatchHandle(defaultErrNo int, extra interface{}) (resp Responser) {
	dispatchEventMessage(msg)
	return
}