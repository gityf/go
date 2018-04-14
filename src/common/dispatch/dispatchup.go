package dispatch

import (
	"time"
	"runtime/debug"
	"common/statistic"
	logger "github.com/xlog4go"
	"common/errorcode"
)

/*
*	负责转发消息：异步，同步
*/

type MessageDispatcher interface {
	//同步分发同步消息，直接与poolhandler(具体的业务逻辑)交互
	Dispatch(defaultErrNo int, extra interface{}) Responser

	//异步分发异步消息, 与dispathcdown交互
	DispatchDown(defaultErrNo int, extra interface{}) Responser

	//异步处理消息, 与dispathcdown交互
	DispatchHandle(defaultErrNo int, extra interface{}) Responser
}

// 调度入口-分发消息等待返回结构
func DispatchMessageWithResponse(msg MessageDispatcher, sync bool, defaultErrNo int) Responser {

	var r Responser
	if sync == true {
		r = msg.Dispatch(defaultErrNo, "")
	} else {
		r = msg.DispatchDown(defaultErrNo, "")
	}

	return r
}

// 调度入口-分发消息不等待返回
func DispatchHandleNoResponse(dispatcher MessageDispatcher, defaultErrNo int, msg *Message) Responser {
	var resp Responser
	var errCode int
	var clientLogId string
	var logId int64

	info := msg.RequestInfo
	tBegin := info["now"].(time.Time)
	clientLogId = msg.ClientLogId
	logId = msg.LogId

	//在异步队列中的耗时
	dispatchLatency := time.Since(tBegin)
	defer func() {
		latency := time.Since(tBegin)
		if resp != nil {
			errCode = resp.ErrCode()
		} else {
			//unlikely
			errCode = -1
		}
		//捕捉panic
		if err := recover(); err != nil {
			errCode = errorcode.ERRNO_PANIC
			logger.Error("[ASYNC] LogId:%d HandleError# recover errno:%d stack:%s", logId, errCode, string(debug.Stack()))
			statistic.IncPanicCount()
		}
		//发送监控数据
		if errCode != 0 {
			logger.Error("[ASYNC] %v [traceid:%v LogId:%d] errno:%d", info["name"], clientLogId, logId, errCode)
		}

		logger.Info("[ASYNC][do=%v logId:%v # dispatchLatency(ms):%.2f latency(ms):%.2f",
			info["name"], logId, dispatchLatency.Seconds()*1000, latency.Seconds()*1000)
	}()

	resp = dispatcher.DispatchHandle(defaultErrNo, "")
	return resp
}

