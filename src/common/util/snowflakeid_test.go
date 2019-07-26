package util

import (
	"testing"
)

func TestSnowFlakeWorkderId1(t *testing.T) {
	sf := NewSnowFlakeWorkerId()
	id, err := sf.NextId()
	t.Logf("id:%v, err:%v\n", id, err)
	count := 20000000
	i := 0
	for i < count {
		id, err = sf.NextId()
		i++
	}
	t.Logf("max-datacenterid:%v", kMaxDataCenterId)
	t.Logf("max-kSequenceMask:%v", kSequenceMask)
}
