package util

import (
	"errors"
	"fmt"
	"sync"
)

/**
 * @brief 分布式id生成类
 * https://segmentfault.com/a/1190000011282426
 * https://github.com/twitter/snowflake/blob/snowflake-2010/src/main/scala/com/twitter/service/snowflake/IdWorker.scala
 *
 * 64bit id: 0000  0000  0000  0000  0000  0000  0000  0000  0000  0000  0000  0000  0000  0000  0000  0000
 *           ||                                                           ||     ||     |  |              |
 *           |└---------------------------时间戳---------------------------┘└中心-┘└机器--┘  └----序列号-----┘
 *           |
 *         不用
 */

type SnowFlakeWorkerId struct {
	workerId      uint64
	datacenterId  uint64
	sequence      uint64
	timestamp     uint64
	lastTimestamp uint64
	mu            *sync.RWMutex
}

// start timestamp ms [2019-01-01 00:00:00]
const kFromEpoch uint64 = 1546272000000

// bits of workder id
const kWorkerIdBits = 5

// bits of data center id
const kDataCenterIdBits = 5

// bits of sequence
const kSequenceBits = 12

// left shift bits of worker id
const kWorkerIdLeftShift = kSequenceBits

// left shift bits of data center id
const kDataCenterIdLeftShift = kWorkerIdLeftShift + kWorkerIdBits

// left shift bits of timestamp
const kTimestampLeftShift = kDataCenterIdLeftShift + kDataCenterIdBits

// max value of data center id is 31
const kMaxDataCenterId = -1 ^ (-1 << kDataCenterIdBits)

// max value of sequence is 4095
const kSequenceMask = -1 ^ (-1 << kSequenceBits)

func NewSnowFlakeWorkerId() *SnowFlakeWorkerId {
	return &SnowFlakeWorkerId{
		workerId:      0,
		datacenterId:  0,
		sequence:      0,
		timestamp:     0,
		lastTimestamp: 0,
	}
}

func (sf *SnowFlakeWorkerId) SetWorkderId(workerId uint64) {
	sf.workerId = workerId
}

func (sf *SnowFlakeWorkerId) SetDataCenterId(datacenterId uint64) {
	sf.datacenterId = datacenterId
}

func (sf *SnowFlakeWorkerId) NextId() (to uint64, err error) {
	//sf.mu.Lock()
	//defer sf.mu.Unlock()

	sf.timestamp = uint64(NowInMs())
	if sf.timestamp < sf.lastTimestamp {
		err = errors.New(fmt.Sprintf("clock is moving backwards.  Rejecting requests until %v", sf.lastTimestamp))
		return
	} else if sf.timestamp == sf.lastTimestamp {
		// same timestamp, generate sequence in one ms.
		sf.sequence = (sf.sequence + 1) & kSequenceMask
		if 0 == sf.sequence {
			sf.timestamp = tilNextMs(sf.lastTimestamp)
		}
	} else {
		sf.sequence = 0
	}
	sf.lastTimestamp = sf.timestamp
	to = ((sf.timestamp - kFromEpoch) << kTimestampLeftShift) | (sf.datacenterId << kDataCenterIdLeftShift) | (sf.workerId << kWorkerIdLeftShift) | sf.sequence
	return
}

func tilNextMs(lastTimestamp uint64) (timestamp uint64) {
	timestamp = uint64(NowInMs())
	for timestamp <= lastTimestamp {
		timestamp = uint64(NowInMs())
	}
	return
}
