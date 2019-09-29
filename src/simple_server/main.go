package main

import (
    "fmt"
    "simple_worker"
)

func run(worker simple_worker.Worker) {
    worker.Run()
}

func main() {
    fmt.Println("start")
    worker := make([]simple_worker.Worker, 2)
    worker[0] = simple_worker.EchoWorker{"hello world"}
    worker[1] = simple_worker.CountWorker{1}
    for _, v := range worker {
        run(v)
    }
}

