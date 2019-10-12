package base

import (
    "fmt"
    "base"
    "protocol"
)

type Server struct {
    rpc map[uint32]*base.RpcMethod
    is_running bool
    quit chan struct{}
    channel *TCPChannel
}

func (svr *Server) RegisterMethods(methods []*base.RpcMethod) error {
    for _, method := range methods {
        cmd := method.Cmd()
        if _, ok := svr.rpc[cmd]; ok {
            return fmt.Errorf("cmd:%X is exit when regist", cmd)
        }
        svr.rpc[cmd] = method
    }
    return nil
}

func (svr *Server) Run() error {
    wg := sync.WaitGroup{}
	wg.Add(1)
	go svr.channel.Serve(&wg, svr.quit, svr.HandleRequest)
	wg.Wait()
    // todo
    return nil
}

func (svr *Server) Stop() error {
    if svr.is_running {
        svr.is_running = false
        close(svr.quit)
    }
    // todo 是否还需要做其他的事情
    return nil
}

func (svr *Server) HandleRequest(req *protocol.Pkg) error {
    return svr.channel.SendPkg(req.Id, req)
}

func (svr *Server) SendResponse(rsp *protocol.Message) error {
    return svr.channel.SendPkg(rsp.Id, rsp)
}

func NewServer(addr string, cache_buf_len uint32) *Server {
    svr := &Server{}
    svr.methods = make(map[uint32]*base.RpcMethod)
    svr.quit = make(chan struct{})
    svr.is_running = false
    svr.channel = NewTCPChannel(addr, cache_buf_len)
    return svr
}


