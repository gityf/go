package gomonitor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type TGoProcStats struct {
	Goroutines    int    `json:"goroutine_num"`
	Mem_Allocated uint64 `json:"mem_allocated"`
	Mem_Objects   uint64 `json:"mem_objects"`
	Mem_Mallocs   uint64 `json:"mem_mallocs"`
	Mem_Heap      uint64 `json:"mem_heap"`
	Mem_Stack     uint64 `json:"mem_stack"`
	Gc_Num        uint32 `json:"gc_num"`
	Gc_Pause      uint64 `json:"gc_pause"`
	Gc_Next       uint64 `json:"gc_next"`
	Cgo           int64  `json:"cgo"`
	Fds           int    `json:"fds"`
	ThreadCreate  int    `json:"threadcreate"`
	Block         int    `json:"block"`
}

var (
	lastNumGc      uint32
	lastPauseTotal uint64
	lastCgoCall    int64
	lastMallocs    uint64
)

func GetRuntimeStats() *TGoProcStats {
	var (
		mem   runtime.MemStats
		stats TGoProcStats
	)

	runtime.ReadMemStats(&mem)

	stats.Mem_Mallocs = mem.Mallocs - lastMallocs
	lastMallocs = mem.Mallocs
	stats.Mem_Allocated = mem.Alloc
	stats.Mem_Objects = mem.HeapObjects
	stats.Mem_Heap = mem.HeapAlloc
	stats.Mem_Stack = mem.StackInuse

	stats.Gc_Num = mem.NumGC - lastNumGc
	lastNumGc = mem.NumGC
	stats.Gc_Pause = mem.PauseTotalNs - lastPauseTotal
	lastPauseTotal = mem.PauseTotalNs
	stats.Gc_Next = mem.NextGC

	stats.Goroutines = runtime.NumGoroutine()
	temp := runtime.NumCgoCall()
	stats.Cgo = temp - lastCgoCall
	lastCgoCall = temp
	stats.Fds = openFileCnt()

	stats.ThreadCreate, _ = runtime.ThreadCreateProfile(nil)
	stats.Block, _ = runtime.BlockProfile(nil)
	return &stats
}

func openFileCnt() int {
	for i := 0; i < 2; i++ {
		out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
		if err == nil {
			return bytes.Count(out, []byte("\n"))
		}
	}
	return 0
}