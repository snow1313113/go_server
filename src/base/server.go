package base

import (
    "fmt"
    "sync"
    "protocol"
    pb "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

func (svr *Server) Registered(service interface{}, service_desc *descriptor.ServiceDescriptorProto) error {
    return collectServiceMethod(service, service_desc, svr.registerhandle)
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
        err := pb.Unmarshal(pkg.Body, req)
        if err != nil {
            fmt.Println(pkg.Head.Cmd, "Unmarshal err: ", err)
            return err
        }
    }

    rsp := method.NewRsp()
    err := method.Call(pkg.Head.Id, req, rsp)

    rsp_pkg := protocol.Pkg{}
    rsp_pkg.Head.Id = pkg.Head.Id
    rsp_pkg.Head.Cmd = pkg.Head.Cmd
    rsp_pkg.Head.Seq = pkg.Head.Seq

    if err != nil {
        fmt.Println("cmd: ", pkg.Head.Cmd, " Call err: ", err)
        // 返回的error必须是RpcError
        rpc_err, ok := err.(*RpcError)
        if !ok {
            fmt.Println("cmd: ", pkg.Head.Cmd, " err:[", err, "] change error")
            return err
        }
        rsp_pkg.Head.Ret = rpc_err.Code
    } else {
        rsp_pkg.Head.Ret = 0
    }

    rsp_pkg.Body, err = pb.Marshal(rsp)
    if err != nil {
        fmt.Println(rsp, " Marshal err: ", err)
        return err
    }
    rsp_pkg.Head.BodyLen = uint32(len(rsp_pkg.Body))

    svr.SendResponse(&rsp_pkg)
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


