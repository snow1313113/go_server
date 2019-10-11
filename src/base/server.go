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
	for {
		select {
		case <-quit:
			return nil
		default:
            // todo 处理请求包，如果以后有跨级通信的话，这里还要处理rpc回调包
		}
	}

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
    // todo
}

func (svr *Server) SendResponse(rsp protocol.Message) error {
    // todo
}

func NewServer() (*Server, error) {
	svr := &Server{}
	svr.methods = make(map[uint32]*base.RpcMethod)
	svr.quit = make(chan struct{})
	svr.is_running = false
	return svr, nil
}


