package main

import (
    "protocol"
)

type ExampleService struct{}

func (s *ExampleService) Hello(id uint32, req *protocol.HelloReq, rsp *protocol.HelloRsp) error {
    log.Debug("handle hello rpc,: %v", req)
    rsp.Name = req.GetName()
    rsp.Num = req.GetNum()
    rsp.Seq = 1
    return nil
}

func (s *ExampleService) Echo(id uint32, req *protocol.EchoReq, rsp *protocol.EchoRsp) error {
    log.Debug("handle echo rpc,: %v", req)
    rsp.RspInfo = req.GetInfo()
    rsp.Seq = 2
    return nil
}


