package common

import "context"

// WorkerPool is an interface for a worker pool taken from "github.com/gammazero/workerpool"
type WorkerPool interface {
	Size() int
	Stop()
	StopWait()
	Stopped() bool
	Submit(task func())
	SubmitWait(task func())
	WaitingQueueSize() int
	Pause(ctx context.Context)
}
