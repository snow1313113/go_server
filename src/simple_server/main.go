package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "base"
    "protocol"
)

type IExampleService interface {
    // HelloReq之类的得实现了protocol.Message接口
    // todo 传参数指针是为了减少拷贝和能把结果带回去，但是实际上是不是把rsp放在返回值回更好
    Hello(*protocol.HelloReq, *protocol.HelloRsp) error
    Hello2(*protocol.HelloReq, *protocol.HelloRsp) error
}

type ExampleService []uint32

func (s *ExampleService) Hello(req *protocol.HelloReq, rsp *protocol.HelloRsp) error {
    fmt.Println("handle hello rpc, req: ", req)
    rsp.Name = req.Name
    rsp.Num = req.Num
    rsp.Seq = 1
    return nil
}

func (s *ExampleService) Hello2(req *protocol.HelloReq, rsp *protocol.HelloRsp) error{
    fmt.Println("handle hello2 rpc, req: ", req)
    rsp.Name = req.Name
    rsp.Num = req.Num
    rsp.Seq = 2
    return nil
}

func main() {
    fmt.Println("start")

    sig_chan := make(chan os.Signal, 1)
    signal.Notify(sig_chan, syscall.SIGINT, syscall.SIGTERM)

    // 这里指定rpc的cmd，好丑的方式，如果用pb描述的话就应该能和rpc写在一起
    example_service := ExampleService{0x01, 0x02}

    svr := base.NewServer("localhost:1234", 512)
    // 注意了，ExampleService的方法是用指针作为reciver的，所以这里要传指针给空接口，
    // 不然反射出来的就不是IExampleService了
    err := svr.Registered(example_service, &example_service)
    if err != nil {
        fmt.Println("regist service err: ", err)
        return
    }
    fmt.Println("registred succeed")

    go svr.Run()

    // 等待收到退出信号
    for sig := range sig_chan {
        fmt.Println("signal : ", sig)
        break
    }

    // 直接调用stop
    svr.Stop()
    fmt.Println("stop svr")
}

