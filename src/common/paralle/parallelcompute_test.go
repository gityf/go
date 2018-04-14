package paralle

import (
	"errors"
	"testing"
	"global"
)

type TestPCReq struct {
	Sid string
	Id  int
}
type TestPCRes struct {
	Err error
	Id  int
}

func PackTestPC(params interface{}, c chan interface{}) (resErr error) {
	var datainfo *TestPCReq
	switch params.(type) {
	case *TestPCReq:
		datainfo = params.(*TestPCReq)
	default:
		return errors.New("params format err")
	}
	res := &TestPCRes{
		Err: nil,
		Id:  datainfo.Id * 2,
	}
	c <- res
	return nil
}
func UnPackTestPC(params []interface{}, dids *[]int) error {
	var one *TestPCRes
	for _, res := range params {
		switch res.(type) {
		case *TestPCRes:
			one = res.(*TestPCRes)
		default:
			return errors.New("params format err")
		}
		if one.Err != nil {
			continue
		}
		*dids = append(*dids, one.Id)
	}
	return nil
}

func Test_ParallelCompute(t *testing.T) {
	global.MainProcIsRunning = 1
	max := 1000
	infos := make([]interface{}, 0, max)
	for i := 0; i < max; i++ {
		info := &TestPCReq{Sid: "test", Id: i}
		//infos[i] = info
		infos = append(infos, info)
	}
	result, err := ParallelCompute(PackTestPC, infos, 500)
	if err != nil {
		t.Error("ParallelCompute errorcode: %v", err)
	}
	dids := make([]int, 0, max)
	err = UnPackTestPC(result, &dids)
	if err != nil {
		t.Error("UnPackTestPC errorcode")
	}
	for i := 0; i < max; i++ {
		if (i * 2) != dids[i] {
			t.Error("Computer errorcode")
		}
	}
}
