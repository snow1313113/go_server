package main

import (
    "fmt"
    "simple_worker"
)

func run(worker simple_worker.Worker) {
    fmt.Println(worker.Run())
}

func main() {
    echo_worker := simple_worker.EchoWorker{"hello world"}
    fmt.Println("start")
    run(echo_worker)
}

