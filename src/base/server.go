package base

import (
    "fmt"
    "sync"
    "protocol"
)

type Server struct {
    rpc map[uint32]*RpcMethod
    is_running bool
    channel *TCPChannel
}

func (svr *Server) registerhandle(methods []*RpcMethod) error {
    for _, method := range methods {
        cmd := method.Cmd()
        if _, ok := svr.rpc[cmd]; ok {
            return fmt.Errorf("cmd:%X is exit when regist", cmd)
        }
        svr.rpc[cmd] = method
    }
    return nil
}

func (svr *Server) Registered(cmds []uint32, service interface{}) error {
    return RegisterService(cmds, service, svr.registerhandle)
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

func (svr *Server) HandleRequest(pkg *protocol.Pkg) error {
    method, ok := svr.rpc[pkg.Head.Cmd]
    if !ok {
        return fmt.Errorf("cmd:%X is not regist", pkg.Head.Cmd)
    }

    req := method.NewReq()
    if pkg.Body != nil {
        err := req.Parse(pkg.Body)
        if err != nil {
            fmt.Println(pkg.Head.Cmd, "parse err: ", err)
            return err
        }
    }

    rsp := method.NewRsp()
    err := method.Call(req, rsp)
    if err != nil {
        fmt.Println(pkg.Head.Cmd, "Call err: ", err)
        return err
    }

    return nil
}

func (svr *Server) SendResponse(rsp *protocol.Pkg) error {
    return svr.channel.SendPkg(rsp.Head.Id, rsp)
}

func NewServer(addr string, cache_buf_len uint32) *Server {
    svr := &Server{}
    svr.rpc = make(map[uint32]*RpcMethod)
    svr.is_running = false
    svr.channel = NewTCPChannel(addr, cache_buf_len)
    return svr
}


