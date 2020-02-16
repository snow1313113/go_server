package base

import (
    "fmt"
    "reflect"
    pb "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/protoc-gen-go/descriptor"
    "protocol"
)

type RpcMethod struct {
    cmd uint32
    name string
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

func (rpc *RpcMethod) NewReq() pb.Message {
    // register的时候判断过了Kind是指针，所以可以直接放心的调用Elem
    return reflect.New(rpc.req_type.Elem()).Interface().(pb.Message)
}

func (rpc *RpcMethod) NewRsp() pb.Message {
    // register的时候判断过了Kind是指针，所以可以直接放心的调用Elem
    return reflect.New(rpc.rsp_type.Elem()).Interface().(pb.Message)
}

// 这里如果传入context，则需要业务逻辑处理context，这样不好，所以暂时干脆不传入了
func (rpc *RpcMethod) Call(gid uint32, req pb.Message, rsp pb.Message) error {
    args := []reflect.Value{rpc.receiver, reflect.ValueOf(gid), reflect.ValueOf(req), reflect.ValueOf(rsp)}
    ret := rpc.method.Func.Call(args)
    if ret == nil {
        return NewRpcError(-1, "Func.Call ret expect []value, but is nil")
    }

    if ret[0].Interface() != nil {
        return NewRpcError(-1, ret[0].Interface().(error).Error())
    }
    return nil
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
        panic(fmt.Sprintln("desc method num %d != impl method num %d", desc_methods_num, st.NumMethod()))
    }
    rpc_methods := make([]*RpcMethod, 0, st.NumMethod())
    for i := 0; i < st.NumMethod(); i++ {
        method := st.Method(i)
        method_type := method.Type

        method_desc, ok := method_desc_map[method.Name]
        if !ok {
            panic(fmt.Sprintln("can not found method %s in desc", method.Name))
        }

        // 所有的method 都有4个参数receiver(实际执行的对象), gid（包头的id), request(请求包), response(返回包)
        if method_type.NumIn() != 4 {
            panic(fmt.Sprintln("method has wrong number of params:", method_type.NumIn()))
        }

        receiver_type := method_type.In(0)
        // 调用对象要求是指针
        if receiver_type.Kind() != reflect.Ptr {
            panic(fmt.Sprintln("method", method.Name, "receiver type not a pointer:", receiver_type))
        }

        gid_type := method_type.In(1)
        // gid参数要求是uint32
        if gid_type.Kind() != reflect.Uint32 {
            panic(fmt.Sprintln("method", method.Name, "gid type not uint32:", gid_type))
        }

        req_type := method_type.In(2)
        // 请求包参数要求是指针
        if req_type.Kind() != reflect.Ptr {
            panic(fmt.Sprintln("method", method.Name, "request type not a pointer:", req_type))
        }

        rsp_type := method_type.In(3)
        // 返回包参数要求是指针
        if rsp_type.Kind() != reflect.Ptr {
            panic(fmt.Sprintln("method", method.Name, "response type not a pointer:", rsp_type))
        }

        option_value, err := pb.GetExtension(method_desc.GetOptions(), protocol.E_CMD)
        if err != nil {
            return err
        }

        cmd := (uint32)(*(option_value.(*uint32)))
        rpc_method := &RpcMethod{
            cmd,
            method.Name,
            reflect.ValueOf(service),
            req_type,
            rsp_type,
            method,
        }

        rpc_methods = append(rpc_methods, rpc_method)
    }

    if len(rpc_methods) != desc_methods_num {
        panic("methods number not match")
    }

    return register(rpc_methods)
}
