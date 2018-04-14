package dispatch

import (
	"global"
	logger "github.com/xlog4go"
)
/**
Step:
 // 1. declare message handler function.
 func msgHandler(eventMsg *EventMessage) Responser {
  ...
 }

 // 2. init dispatch
 initDispatch()
 // 3. get dispatcher
 ddws = GetDefaultDispatcher()
 // 4. new message and send to dispatcher queue.
 eventMsg := &EventMessage{
	CallFunc: msgHandler,
 }
 ddws.SendMessage(eventMsg)
*/
// to manage workspace of multi-biz
var ddwsMgr *DispatchDownWorkspaceManager

func initDispatch(eventConcurrency, eventMessageQueueLen int) (err error) {
	//init DispatchDownWorkspaceManager
	ddwsMgr = NewDispatchDownWorkspaceManager()

	//Dispatch DownStream initial Begin...
	var ddws *DispatchDownWorkspace
	queueConNumInfo := &QueueConNumInfo{
		eventConCurrency:              eventConcurrency,
		eventMessgageQueueLen:         eventMessageQueueLen,
	}
	ddws, err = NewDispatchDownWorkspace(queueConNumInfo)
	if err != nil {
		logger.Error("init DispatchDownWorkspace  Failed:%s", err.Error())
		return
	}
	ddws.Start()
	ddwsMgr.Insert(global.DEFAULT_SID, ddws)
	return
}

func GetDefaultDispatcher() (ddws *DispatchDownWorkspace) {
	ddws, _ = ddwsMgr.GetDispatchDownWorkspace(global.DEFAULT_SID)
	return
}