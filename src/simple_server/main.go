package main

import (
    "fmt"
    "sync"
    "base"
    "protocol"
)

type IExampleService interface {
    // HelloReq之类的得实现了protocol.Message接口
    // todo 传参数指针是为了减少拷贝和能把结果带回去，但是实际上是不是把rsp放在返回值回更好
    Hello(*HelloReq, *HelloRsp) error
    Hello2(*Hello2Req, *Hello2Rsp) error
}

type ExampleService []uint32

func (s *ExampleService) Hello(*HelloReq, *HelloRsp) error {
    // todo
    return nil
}

func (s *ExampleService) Hello2(*Hello2Req, *Hello2Rsp) error{
    // todo
    return nil
}

func main() {
    fmt.Println("start")

    // 这里指定rpc的cmd，好丑的方式，如果用pb描述的话就应该能和rpc写在一起
    example_service := ExampleService{0x01, 0x02}

    svr := NewServer("localhost:1234", 512)
    err := svr.RegisterService(example_service, example_service)
    if err != nil {
        fmt.Println("regist service err: ", err)
        return
    }

    // todo 这样处理是不行的，而且还有stop没有处理，还是得起一个deamon进程去跑，deamon进程的方式还没研究
    wg := sync.WaitGroup
    wg.Add(1)
    go svr.Run()
    wg.Wait()
}

