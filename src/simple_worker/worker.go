package simple_worker

import (
    "fmt"
)

type Worker interface {
    Init() int
    Run() int
    Stop() bool
}

type EchoWorker struct {
    Word string
}

func (w EchoWorker) Init() int {
    return 0
}

func (w EchoWorker) Run() int {
    fmt.Println("EchoWorker", w.Word)
    return 0
}

func (w EchoWorker) Stop() bool {
    return true
}

type CountWorker struct {
    ReqCount int
}

func (w CountWorker) Init() int {
    return 0
}

func (w CountWorker) Run() int {
    w.ReqCount++
    fmt.Println("CountWorker", w.ReqCount)
    return 0
}

func (w CountWorker) Stop() bool {
    return true
}

