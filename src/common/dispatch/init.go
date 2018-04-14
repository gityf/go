package dispatch

import (
	"global"
	logger "github.com/xlog4go"
)

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