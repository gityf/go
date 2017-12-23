package statistic

import (
	"sync/atomic"
	"time"
)

type Statistic struct {
	Id int
	// qps of request
	ReqCount int64
	EnlargeCnt int64
	JobCount int64
	QPSTBegin int64
	PanicCount int64
}

func NewStatistic(id int) *Statistic {
	var ss Statistic
	ss.ReqCount = 0
	ss.EnlargeCnt = 0
	ss.JobCount = 0
	ss.Id = id
	ss.QPSTBegin = 0
	ss.PanicCount = 0
	return &ss
}

func (ss *Statistic) GetId() int {
	return ss.Id
}

func (ss *Statistic) GetReqCount() int64 {
	return ss.ReqCount
}

func (ss *Statistic) GetEnlargeCnt() int64 {
	return ss.EnlargeCnt
}

func (ss *Statistic) GetJobCount() int64 {
	return ss.JobCount
}

func (ss *Statistic) ReInit() {
	ss.ReqCount = 0
	ss.EnlargeCnt = 0
	ss.JobCount = 0
	ss.QPSTBegin = 0
}

var StatInfo []*Statistic
var StatInfoPointer atomic.Value

func init() {
	StatInfo = append(StatInfo, NewStatistic(0))
	StatInfo = append(StatInfo, NewStatistic(1))

	StatInfoPointer.Store(StatInfo[0])
}

//访问次数
func (ss *Statistic) IncReqCount() {
	atomic.AddInt64(&ss.ReqCount, 1)
}

func (ss *Statistic) IncEnlargeCnt() {
	atomic.AddInt64(&ss.EnlargeCnt, 1)
}

func (ss *Statistic) IncJobCount() {
	atomic.AddInt64(&ss.JobCount, 1)
}

func (ss *Statistic) IncPanicCount() {
	atomic.AddInt64(&ss.PanicCount, 1)
}

func IncrReqCount(){
	StatInfoPointer.Load().(*Statistic).IncReqCount()
}

func IncPanicCount(){
	StatInfoPointer.Load().(*Statistic).IncPanicCount()
}

func GetReqQpsAndReset() int64 {
	id := StatInfoPointer.Load().(*Statistic).GetId()
	nowTime := time.Now().UnixNano() / int64(time.Millisecond)
	StatInfo[1-id].ReInit()
	StatInfo[1-id].QPSTBegin = nowTime
	//指向另外一个Info
	StatInfoPointer.Store(StatInfo[1-id])

	ssPtr := StatInfo[id]
	return (ssPtr.GetReqCount() * 1000) / (nowTime - ssPtr.QPSTBegin)
}

func GetReqQps() int64 {
	return (StatInfoPointer.Load().(*Statistic).GetReqCount() * 1000) /
		(time.Now().UnixNano() / int64(time.Millisecond) - StatInfoPointer.Load().(*Statistic).QPSTBegin)
}

func GetJobCount() int64 {
	return StatInfoPointer.Load().(*Statistic).GetJobCount()
}

func GetEnlargeCnt() int64 {
	return StatInfoPointer.Load().(*Statistic).GetEnlargeCnt()
}

func GetStatistic () *Statistic {
	return StatInfoPointer.Load().(*Statistic)
}