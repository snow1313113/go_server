package main

import (
    "fmt"
    "net"
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

    _, err = conn.Write([]byte("say hello!"))
    if err != nil {
        fmt.Println("write err : ", err)
        return
    }
}
