package main

import (
    "fmt"
    "net"
    "syscall"
    "protocol/pkg"
)

func handleMsg(buf []byte, length int){
    fmt.Println("handle msg: ", string(buf[:length]))
}

func dealConn(conn net.Conn) {
    fmt.Println("new connect ...from: ", conn.RemoteAddr().String())
    defer conn.Close()
    for {
        buf := make([]byte, 512)
        buf_len, err := conn.Read(buf)
        switch err {
        case nil:
            handleMsg(buf, buf_len)
        case syscall.EAGAIN:
            // 需要把收到的数据缓存拼接一下
            continue
        default:
            fmt.Println("read err : ", err)
            return
        }
    }
}

func main() {
    fmt.Println("start server ...")
    listener, err := net.Listen("tcp", "localhost:1234")
    if err != nil {
        fmt.Println("listen err : ", err)
        return
    }
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("accept err : ", err)
        }
        go dealConn(conn)
    }
}
