package paralle

import (
	logger "github.com/xlog4go"
	"errors"
	"runtime/debug"
	"time"
	"global"
)

const (
	MAX_PARALLEL_GORUNTINE = 100
	ERRNO_PANIC = 500
)

type CallFuncName func(params interface{}, c chan interface{}) error

func ParallelCompute(callfun CallFuncName, infos []interface{}, timeout_ms int) ([]interface{}, error) {
	var err error
	length := len(infos)
	result := make([]interface{}, 0, length)

	max := length / MAX_PARALLEL_GORUNTINE
	mod := length % MAX_PARALLEL_GORUNTINE
	if mod > 0 {
		max++
	}
	for index := 0; index < max; index++ {
		beg := index * MAX_PARALLEL_GORUNTINE
		end := (index + 1) * MAX_PARALLEL_GORUNTINE
		if end > length {
			end = length
		}
		max_num := end - beg
		retch := make(chan interface{}, max_num)
		for tmp := beg; tmp < end; tmp++ {
			if global.MainProcIsRunning != 1 {
				return result, errors.New("received exit signal")
			}

			time.Sleep(1 * time.Millisecond)

			go func(tmp int) {
				defer func() {
					if err := recover(); err != nil {
						logger.Error("callfun  unknown errorcode,errno:%d errmsg:%v, stack:%s", ERRNO_PANIC, err, string(debug.Stack()))
					}
				}()

				callfun(infos[tmp], retch)
			}(tmp)
		}

		num := 0
		for {
			select {
			case retdata := <-retch:
				num++
				result = append(result, retdata)
				//case <-timeout:
			case <-time.After(time.Duration(timeout_ms) * time.Millisecond):
				if num < max_num {
					logger.Error("ParallelCompute timeout")
					err = errors.New("ParallelCompute timeout")
					return result, err
				}
			}
			if num == max_num {
				break
			}
			if global.MainProcIsRunning != 1 {
				return result, errors.New("received exit signal")
			}
		}
	}
	return result, nil
}