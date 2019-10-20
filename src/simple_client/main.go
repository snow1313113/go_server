package main

import (
    "fmt"
    "net"
    "time"
    "protocol"
)

func main() {
    fmt.Println("start client ...")
    conn, err := net.Dial("tcp", ":1234")
    if err != nil {
        fmt.Println("connect err : ", err)
        return
    }

    defer func() {
        fmt.Println("client finish")
        conn.Close()
    }()

    for i := 1; i < 10; i++ {
        pkg := protocol.Pkg{}

        pkg.Head.Id = 0
        pkg.Head.Cmd = 0x01
        pkg.Head.Seq = uint32(i)

        req := protocol.HelloReq{uint32(i*100), 10}
        pkg.Body, err = req.Bytes()
        if err != nil {
            fmt.Println("req bytes err:", err)
            return
        }

        pkg.Head.BodyLen = uint32(len(pkg.Body))

        b, err := pkg.Bytes()
        if err != nil {
            fmt.Println("pkg bytes err:", err)
            return
        }

        fmt.Println("here")
        _, err = conn.Write(b)
        if err != nil {
            fmt.Println("write err : ", err)
            return
        }

        time.Sleep(time.Duration(3)*time.Second)
    }
}
