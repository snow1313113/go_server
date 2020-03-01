package base

import (
    "sync"
    "protocol"
    "utils"
    pb "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var (
    log *utils.Logger
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
            log.Error("duplicate cmd[0x%08x]", cmd)
            return NewRpcError(-1, "")
        }
        svr.rpc[cmd] = method
        log.Debug("regist cmd[0x%08x]", cmd)
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
    log.Debug("recv: gid[%d] cmd[0x%08x] bytelen[%d]", pkg.Head.Id, pkg.Head.Cmd, len(pkg.Body))
    method, ok := svr.rpc[pkg.Head.Cmd]
    if !ok {
        log.Error("cmd[0x%08x] is not found", pkg.Head.Cmd)
        return NewRpcError(-1, "")
    }

    req := method.NewReq()
    if pkg.Body != nil {
        err := pb.Unmarshal(pkg.Body, req)
        if err != nil {
            log.Error("cmd 0x%x Unmarshal err: %v", pkg.Head.Cmd, err)
            return err
        }
    }

    rsp := method.NewRsp()
    ret_code, err := method.Call(pkg.Head.Id, req, rsp)
    if err != nil {
        log.Error("cmd 0x%x Call err: %v", pkg.Head.Cmd, err)
        return err
    }

    if method.NeedRsp() {
        rsp_pkg := protocol.Pkg{}
        rsp_pkg.Head.Id = pkg.Head.Id
        rsp_pkg.Head.Cmd = pkg.Head.Cmd
        rsp_pkg.Head.Seq = pkg.Head.Seq
        rsp_pkg.Head.Ret = ret_code
        if ret_code == 0 {
            rsp_pkg.Body, err = pb.Marshal(rsp)
            if err != nil {
                log.Error("cmd 0x%x Marshal err: %v", pkg.Head.Cmd, err)
                return err
            }
        }
        rsp_pkg.Head.BodyLen = uint32(len(rsp_pkg.Body))

        svr.sendResponse(&rsp_pkg)
    }

    return nil
}

func (svr *Server) sendResponse(rsp *protocol.Pkg) error {
    return svr.channel.SendPkg(rsp.Head.Id, rsp)
}

func (svr *Server) SendMsg(gid uint32, cmd uint32, msg pb.Message) error {
    pkg := protocol.Pkg{}
    pkg.Head.Id = gid
    pkg.Head.Cmd = cmd
    pkg.Head.Seq = 0
    pkg.Head.Ret = 0
    body, err := pb.Marshal(msg)
    if err != nil {
        log.Error("cmd 0x%x Marshal err: %v", pkg.Head.Cmd, err)
        return err
    }
    pkg.Body = body
    pkg.Head.BodyLen = uint32(len(pkg.Body))
    return svr.channel.SendPkg(gid, &pkg)
}

func NewServer(addr string, cache_buf_len uint32, logger *utils.Logger) *Server {
    // todo 在这里初始化全局的变量，不是很好
    log = logger
    if log == nil {
        return nil
    }

    svr := &Server{}
    svr.rpc = make(map[uint32]*RpcMethod)
    svr.is_running = false
    svr.channel = NewTCPChannel(addr, cache_buf_len)
    return svr
}


