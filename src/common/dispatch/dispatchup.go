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
func DispatchHandleNoResponse(dispatcher MessageDispatcher, defaultErrNo int, tBegin time.Time) Responser {
	var resp Responser
	var errCode int
	var logId int64

	//在异步队列中的耗时
	dispatchLatency := time.Since(tBegin)
	defer func() {
		latency := time.Since(tBegin)
		//捕捉panic
		if err := recover(); err != nil {
			errCode = errorcode.ERRNO_PANIC
			logger.Error("[ASYNC] LogId:%d HandleError# recover errno:%d stack:%s", logId, errCode, string(debug.Stack()))
			statistic.IncPanicCount()
		}

		logger.Info("[ASYNC][logId:%v # dispatchLatency(ms):%.2f latency(ms):%.2f",
			logId, dispatchLatency.Seconds()*1000, latency.Seconds()*1000)
	}()

	resp = dispatcher.DispatchHandle(defaultErrNo, "")
	return resp
}


