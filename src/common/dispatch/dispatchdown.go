package dispatch


import (
	"fmt"
	"runtime/debug"
	"sync"
	"common/statistic"
	logger "github.com/xlog4go"
	"common/errorcode"
)

type DispatchDownWorkspaceManager struct {
	ddwsMap map[string]*DispatchDownWorkspace
}

func NewDispatchDownWorkspaceManager() *DispatchDownWorkspaceManager {

	var ddwsMgr DispatchDownWorkspaceManager
	ddwsMgr.ddwsMap = make(map[string]*DispatchDownWorkspace)

	return &ddwsMgr
}

func (ddwsMgr *DispatchDownWorkspaceManager) GetDispatchDownWorkspace(sid string) (*DispatchDownWorkspace, error) {

	ddws, has := ddwsMgr.ddwsMap[sid]
	if has == true {
		return ddws, nil
	}
	//TODO
	return nil, nil
}

func (ddwsMgr *DispatchDownWorkspaceManager) Insert(sid string, ddws *DispatchDownWorkspace) {
	ddwsMgr.ddwsMap[sid] = ddws
}

func (ddwsMgr *DispatchDownWorkspaceManager) Close() {
	for _, ddws := range ddwsMgr.ddwsMap {
		ddws.Close()
		<-ddws.GetCloseCompleteSignal()
	}

}

type DispatchDownWorkspace struct {
	eventMessgageQueue    chan *EventMessage
	eventMessgageQueueLen int
	eventConCurrency      int

	closeSignal         chan interface{}
	closeCompleteSignal chan interface{}
	wg                  sync.WaitGroup
}

type QueueConNumInfo struct {
	//事件队列大小和协程数
	eventMessgageQueueLen int
	eventConCurrency      int
}

func NewDispatchDownWorkspace(qcnInfo *QueueConNumInfo) (*DispatchDownWorkspace, error) {
	var ddws DispatchDownWorkspace

	ddws.eventConCurrency = qcnInfo.eventConCurrency
	ddws.eventMessgageQueueLen = qcnInfo.eventMessgageQueueLen
	ddws.eventMessgageQueue = make(chan *EventMessage, ddws.eventMessgageQueueLen)

	ddws.closeSignal = make(chan interface{})
	ddws.closeCompleteSignal = make(chan interface{})

	return &ddws, nil
}

func (ddws *DispatchDownWorkspace) Start() {

	for i := 0; i < ddws.eventConCurrency; i += 1 {
		ddws.wg.Add(1)
		go ddws.eventWorker(i)
	}

	logger.Info("%d Worker Goroutine Success Started", ddws.eventConCurrency)
}

func (ddws *DispatchDownWorkspace) eventWorker(idx int) {
	defer ddws.wg.Done()
	logger.Info("DispatchDownEventWorker:%d start process", idx)
	defer logger.Info("DispatchDownEventWorker:%d stop process", idx)

	defer func() {
		if err := recover(); err != nil {
			logger.Error("dispatchdown event unknown errorcode, errno:%d errmsg:%v, stack:%s",
				errorcode.ERRNO_PANIC, err, string(debug.Stack()))
			statistic.IncPanicCount()
		}
	}()

PROCESS_LOOP:
	for {
		select {
		case msg := <-ddws.eventMessgageQueue:
			//dispatchEventMessage(msg)
			DispatchHandleNoResponse(msg, 0, msg.Message)
		case <-ddws.closeSignal:
			break PROCESS_LOOP
		}
	}
}

func (ddws *DispatchDownWorkspace) SendMessage(msg interface{}) (err error) {

	switch msg.(type) {
	case *EventMessage:
		if len(ddws.eventMessgageQueue) > (ddws.eventMessgageQueueLen * 4 / 5) {
			logger.Warn("the ddws.eventMessgageQueue's idle is lower than 20%% sid:%v used:%d,total:%d",
				msg.(*EventMessage).Sid,
				len(ddws.eventMessgageQueue),
				ddws.eventMessgageQueueLen)
		}
		select {
		case ddws.eventMessgageQueue <- msg.(*EventMessage):
		default:
			logger.Fatal("eventMessgageQueue is full")
			err = fmt.Errorf("eventMessgageQueue is full")
		}
	default:
	}
	return
}

func (ddws *DispatchDownWorkspace) Close() {
	close(ddws.closeSignal)
	go func() {
		ddws.wg.Wait()
		close(ddws.closeCompleteSignal)
	}()
}

func (ddws *DispatchDownWorkspace) GetCloseCompleteSignal() <-chan interface{} {
	return ddws.closeCompleteSignal
}