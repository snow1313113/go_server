package main

import (
    "protocol"
)

type ExampleService struct{}

func (s *ExampleService) Hello(id uint32, req *protocol.HelloReq, rsp *protocol.HelloRsp) int32 {
    log.Debug("handle hello rpc,: %v", req)
    rsp.Name = req.GetName()
    rsp.Num = req.GetNum()
    rsp.Seq = 1
    return 0
}

func (s *ExampleService) Echo(id uint32, req *protocol.EchoReq, rsp *protocol.EchoRsp) int32 {
    log.Debug("handle echo rpc,: %v", req)
    rsp.RspInfo = req.GetInfo()
    rsp.Seq = 2
    return 0
}


