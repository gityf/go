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

 // 2. init dispatch manager and add dispatchor. 
 InitDispatchMananger()
 AddDispatcher()

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

func InitDispatchMananger() {
	//init DispatchDownWorkspaceManager
	ddwsMgr = NewDispatchDownWorkspaceManager()
}

func AddDispatcher(eventConcurrency, eventMessageQueueLen int, sid string) (err error) {
	if ddwsMgr == nil {
		InitDispatchMananger()
	}

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
	ddwsMgr.Insert(sid, ddws)
	return
}

func GetDefaultDispatcher() (ddws *DispatchDownWorkspace) {
	ddws, _ = ddwsMgr.GetDispatchDownWorkspace(global.DEFAULT_SID)
	return
}

func GetDispatcherBySid(sid string) (ddws *DispatchDownWorkspace) {
	ddws, _ = ddwsMgr.GetDispatchDownWorkspace(sid)
	return
}

