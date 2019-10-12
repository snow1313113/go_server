package main

import (
    "fmt"
    "base"
    "protocol"
)

type IExampleService interface {
    // HelloReq之类的得实现了protocol.Message接口
    // todo 传参数指针是为了减少拷贝和能把结果带回去，但是实际上是不是把rsp放在返回值回更好
    Hello(*HelloReq, *HelloRsp) error
    Hello2(*Hello2Req, *Hello2Rsp) error
}

type ExampleService struct{}

func (s *ExampleService) Hello(*HelloReq, *HelloRsp) error {
    // todo
    return nil
}

func (s *ExampleService) Hello2(*Hello2Req, *Hello2Rsp) error{
    // todo
    return nil
}

var (
    registed_methods MethodsRegister
)

func main() {
    fmt.Println("start")

    example_service := ExampleService{}
    err := base.RegistedMethods(example_service, registed_methods)
    if err != nil {
        fmt.Println("regist service err: ", err)
        return
    }
}

