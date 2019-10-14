package base

import (
    "fmt"
    "base"
    "protocol"
)

type Server struct {
    rpc map[uint32]*base.RpcMethod
    is_running bool
    channel *TCPChannel
}

func (svr *Server) register(methods []*base.RpcMethod) error {
    for _, method := range methods {
        cmd := method.Cmd()
        if _, ok := svr.rpc[cmd]; ok {
            return fmt.Errorf("cmd:%X is exit when regist", cmd)
        }
        svr.rpc[cmd] = method
    }
    return nil
}

func (svr *Server) RegisterService(cmds []uint32, service interface{}) error {
    return RegisterService(cmds, service, svr.register)
}

func (svr *Server) Run() error {
    svr.is_running = true
    wg := sync.WaitGroup{}
    wg.Add(1)
    go svr.channel.Start(&wg, svr.HandleRequest)
    wg.Wait()
    svr.is_running = false
    return nil
}

func (svr *Server) Stop() error {
    if svr.is_running {
        svr.is_running = false
        svr.channel.ShutDown()
    }
    return nil
}

func (svr *Server) HandleRequest(req *protocol.Pkg) error {
    return svr.channel.SendPkg(req.Id, req)
}

func (svr *Server) SendResponse(rsp *protocol.Pkg) error {
    return svr.channel.SendPkg(rsp.Id, rsp)
}

func NewServer(addr string, cache_buf_len uint32) *Server {
    svr := &Server{}
    svr.methods = make(map[uint32]*base.RpcMethod)
    svr.is_running = false
    svr.channel = NewTCPChannel(addr, cache_buf_len)
    return svr
}


