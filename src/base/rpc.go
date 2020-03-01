package base

import (
    "reflect"
    pb "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/protoc-gen-go/descriptor"
    "protocol"
)

type RpcMethod struct {
    cmd uint32
    name string
    need_rsp bool
    // 最终调用的对象、请求包、返回包都是类型接口，后面通过反射再获取实际类型
    receiver reflect.Value
    req_type reflect.Type
    rsp_type reflect.Type
    method reflect.Method
}

func (rpc *RpcMethod) Cmd() uint32 {
    return rpc.cmd
}

func (rpc *RpcMethod) Name() string {
    return rpc.name
}

func (rpc *RpcMethod) NeedRsp() bool {
    return rpc.need_rsp
}

func (rpc *RpcMethod) NewReq() pb.Message {
    // register的时候判断过了Kind是指针，所以可以直接放心的调用Elem
    return reflect.New(rpc.req_type.Elem()).Interface().(pb.Message)
}

func (rpc *RpcMethod) NewRsp() pb.Message {
    // register的时候判断过了Kind是指针，所以可以直接放心的调用Elem
    return reflect.New(rpc.rsp_type.Elem()).Interface().(pb.Message)
}

// 这里如果传入context，则需要业务逻辑处理context，这样不好，所以暂时干脆不传入了
func (rpc *RpcMethod) Call(gid uint32, req pb.Message, rsp pb.Message) (int32, error) {
    log.Debug("call cmd[0x%08x]", rpc.cmd)
    args := []reflect.Value{rpc.receiver, reflect.ValueOf(gid), reflect.ValueOf(req), reflect.ValueOf(rsp)}
    ret := rpc.method.Func.Call(args)
    if ret == nil {
        return -1, NewRpcError(-1, "Func.Call ret expect []value, but is nil")
    }

    // 先判断一下类型对不对吧，因为Int()的内部实现如果转失败会直接panic
    if ret[0].Kind() != reflect.Int32 {
        return -1, NewRpcError(-1, "method return must be int32")
    }

    return int32(ret[0].Int()), nil
}

// 既然rpc method的结构在这里定义了，那如何解析注册也得在这里了
func collectServiceMethod(service interface{}, service_desc *descriptor.ServiceDescriptorProto, register func([]*RpcMethod) error) error {
    // 收集service desc里面的method，用于后面校验和判断
    desc_methods_num := len(service_desc.GetMethod())
    method_desc_map := make(map[string]*descriptor.MethodDescriptorProto, desc_methods_num)
    for _, method_desc := range service_desc.GetMethod() {
        method_desc_map[method_desc.GetName()] = method_desc
    }

    // 获取service的类型信息，可以从里面提取出所有的method
    st := reflect.TypeOf(service)
    if desc_methods_num != st.NumMethod() {
        log.Error("desc method num %d != impl method num %d", desc_methods_num, st.NumMethod())
        return NewRpcError(-1, "")
    }
    rpc_methods := make([]*RpcMethod, 0, st.NumMethod())
    for i := 0; i < st.NumMethod(); i++ {
        method := st.Method(i)
        method_type := method.Type

        method_desc, ok := method_desc_map[method.Name]
        if !ok {
            log.Error("can not found method %s in desc", method.Name)
            return NewRpcError(-1, "")
        }

        // 所有的method 都有4个参数receiver(实际执行的对象), gid（包头的id), request(请求包), response(返回包)
        if method_type.NumIn() != 4 {
            log.Error("method has wrong number of params : %d", method_type.NumIn())
            return NewRpcError(-1, "")
        }

        receiver_type := method_type.In(0)
        // 调用对象要求是指针
        if receiver_type.Kind() != reflect.Ptr {
            log.Error("method(%s) receiver type not a pointer", method.Name, receiver_type)
            return NewRpcError(-1, "")
        }

        gid_type := method_type.In(1)
        // gid参数要求是uint32
        if gid_type.Kind() != reflect.Uint32 {
            log.Error("method(%s) gid type not uint32", method.Name, gid_type)
            return NewRpcError(-1, "")
        }

        req_type := method_type.In(2)
        // 请求包参数要求是指针
        if req_type.Kind() != reflect.Ptr {
            log.Error("method(%s) request type not a pointer", method.Name, req_type)
            return NewRpcError(-1, "")
        }

        rsp_type := method_type.In(3)
        // 返回包参数要求是指针
        if rsp_type.Kind() != reflect.Ptr {
            log.Error("method(%s) response type not a pointer", method.Name, rsp_type)
            return NewRpcError(-1, "")
        }

        option_value, err := pb.GetExtension(method_desc.GetOptions(), protocol.E_CMD)
        if err != nil {
            log.Error("protobuf get extension err: %v", err)
            return NewRpcError(-1, "")
        }

        dont_rsp := false
        dont_rsp_option, err := pb.GetExtension(method_desc.GetOptions(), protocol.E_DONT_RSP)
        if err == nil {
            log.Debug("method(%s) do not need response", method.Name)
            dont_rsp = *(dont_rsp_option.(*bool))
        }

        cmd := (uint32)(*(option_value.(*uint32)))
        rpc_method := &RpcMethod{
            cmd,
            method.Name,
            !dont_rsp,
            reflect.ValueOf(service),
            req_type,
            rsp_type,
            method,
        }

        rpc_methods = append(rpc_methods, rpc_method)
    }

    if len(rpc_methods) != desc_methods_num {
        log.Error("methods number %d not match %d", len(rpc_methods), desc_methods_num)
        return NewRpcError(-1, "")
    }

    return register(rpc_methods)
}
