package sys

import (
	"bytes"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

// GetGoroutineID 获取当前 goroutine 的 ID（非官方方式）
func GetGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := bytes.Fields(buf[:n])
	if len(fields) < 2 {
		return -1
	}
	id, err := strconv.ParseInt(string(fields[1]), 10, 64)
	if err != nil {
		return -1
	}
	return id
}

// GetThreadID 获取当前线程的 OS 线程 ID
func GetThreadID() int64 {
	tid, _, _ := syscall.RawSyscall(syscall.SYS_GETTID, 0, 0, 0)
	return int64(tid)
}

// GetProcessID ...
func GetProcessID() int64 {
	return int64(os.Getpid())
}
