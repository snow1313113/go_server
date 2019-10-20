package base

import (
    "fmt"
    "reflect"
    "errors"
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

func (rpc *RpcMethod) NewReq() protocol.Message {
    // todo 如果是proto的话返回的是proto.Message，这个应该是interface
    // register的时候判断过了Kind是指针，所以可以直接放心的调用Elem
    return reflect.New(rpc.req_type.Elem()).Interface().(protocol.Message)
}

func (rpc *RpcMethod) NewRsp() protocol.Message {
    // todo 如果是proto的话返回的是proto.Message，这个应该是interface
    // register的时候判断过了Kind是指针，所以可以直接放心的调用Elem
    return reflect.New(rpc.rsp_type.Elem()).Interface().(protocol.Message)
}

// 这里如果传入context，则需要业务逻辑处理context，这样不好，所以暂时干脆不传入了
func (rpc *RpcMethod) Call(req protocol.Message, rsp protocol.Message) error {
    args := []reflect.Value{rpc.receiver, reflect.ValueOf(req), reflect.ValueOf(rsp)}
    ret := rpc.method.Func.Call(args)
    if ret == nil {
        return errors.New("Func.Call ret expect []value, but is nil")
    }

    if ret[0].Interface() != nil {
        return ret[0].Interface().(error)
    }
    return nil
}

// 既然rpc method的结构在这里定义了，那如何解析注册也得在这里了
func RegisterService(cmds []uint32, service interface{}, register func([]*RpcMethod) error) error {
    // 获取service的类型信息，可以从里面提取出所有的method
    st := reflect.TypeOf(service)
    if len(cmds) > st.NumMethod(){
        panic(fmt.Sprintln("cmd num > method num", len(cmds), st.NumMethod()))
    }
    rpc_methods := make([]*RpcMethod, 0, st.NumMethod())
    for i := 0; i < st.NumMethod(); i++ {
        method := st.Method(i)
        method_type := method.Type

        // 所有的method 都有三个参数receiver(实际执行的对象), request(请求包), response(返回包)
        if method_type.NumIn() != 3 {
            panic(fmt.Sprintln("method has wrong number of params:", method_type.NumIn()))
        }

        receiver_type := method_type.In(0)
        // 调用对象要求是指针
        if receiver_type.Kind() != reflect.Ptr {
            panic(fmt.Sprintln("method", method.Name, "receiver type not a pointer:", receiver_type))
        }

        req_type := method_type.In(1)
        // 请求包参数要求是指针
        if req_type.Kind() != reflect.Ptr {
            panic(fmt.Sprintln("method", method.Name, "request type not a pointer:", req_type))
        }

        rsp_type := method_type.In(2)
        // 返回包参数要求是指针
        if rsp_type.Kind() != reflect.Ptr {
            panic(fmt.Sprintln("method", method.Name, "response type not a pointer:", rsp_type))
        }

        rpc_method := &RpcMethod{
            cmds[i],
            "todo string",
            reflect.ValueOf(service),
            req_type,
            rsp_type,
            method,
        }

        rpc_methods = append(rpc_methods, rpc_method)
    }

    return register(rpc_methods)
}
