package base

import (
    "fmt"
    "sync"
    "base"
    "protocol"
)

type TCPChannel struct {
    ip_and_port string
    max_buf_len uint32
    // todo 唯一ID怎么定义，简单点的就递增序列分配吧
    generate_id uint32
    // 一个connect有一个缓冲，用来存没有收完的包
    pending_pkg map[uint32][]byte
    // 一个唯一id对应一个connect
    conn_map map[uint32]net.Conn
}

func NewTCPChannel(addr string, buf_len uint32) *TCPChannel {
    channel := &TCPChannel{ip_and_port : addr, max_buf_len : buf_len}
    return channel
}

func (c *TCPChannel) dealConn(conn net.Conn, wg *sync.WaitGroup, quit <-chan struct{}, handler func(*protocol.Pkg)error) {
    // 退出的时候要关闭链接和通知调用协程
    defer func () {
        conn.Close()
        wg.Done()
    }()

    c.generate_id++
    new_id := c.generate_id
    c.conn_map[new_id] = conn
    defer delete(c.conn_map, new_id)

    // todo 还没有日志系统，只能这样打了
    fmt.Println("new connect ...from: ", conn.RemoteAddr().String())
    for {
        select {
        case <-quit:
            // 收到信号就退出
            return
        default:
            pkg := &protocol.Pkg{}
            err := c.RecvPkg(new_id, pkg)
            if err != nil {
                // todo 还没有日志系统，只能这样打了
                fmt.Println("recv pkg err: ", err)
                continue
            }

            // 简单点处理，这里把生成的唯一id给填上
            pkg.Id = new_id
            err = handler(pkg)
            if err != nil {
                fmt.Println("handle pkg err: ", err)
                continue
            }
        }
    }
}

func (c *TCPChannel) Start(wg *sync.WaitGroup, quit <-chan struct{}, handler func(*protocol.Pkg)error) {
    // 如果是协程方式调用，这里可以保证通知到调用的协程这个函数执行完了
    defer wg.Done()

    // 监听tcp端口
    listener, err := net.Listen("tcp", ip_and_port)
    if err != nil {
        fmt.Println("listen err : ", err)
        return
    }

    channel_wg := sync.WaitGroup
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("accept err : ", err)
        }
        channel_wg.Add(1)
        // todo 有新链接就启动一个协程处理，这样处理合理么？
        go dealConn(conn, channel_wg, quit, handler)
    }
    // 等待所有处理单独链接的协程退出
    channel_wg.Wait()
}

func (c *TCPChannel) SendPkg(id uint32, pkg *protocol.Pkg) error {
    if conn, is_exist := c.conn_map[id]; is_exist != nil {
        fmt.Println(id, " not exist")
        return is_exist
    }

    byte_buf, err := pkg.Bytes()
    if err != nil {
        fmt.Println("send err: ", err)
        return err
    }

    _, err = conn.Write(byte_buf)
    if err != nil {
        fmt.Println("write err : ", err)
        return err
    }
    return nil
}

func (c *TCPChannel) RecvPkg(id uint32, pkg *protocol.Pkg) error {
    if conn, is_exist := c.conn_map[id]; is_exist != nil {
        fmt.Println(id, " not exist")
        return is_exist
    }

    // todo buffer可能要固定，不要每次recv都生成一个
    buf := make([]byte, c.max_buf_len)
    buf_len, err := conn.Read(buf)
    switch err {
    case nil:
        err = pkg.Parse(buf[:buf_len])
        if err != nil {
            // todo 还没有日志系统，只能这样打了
            fmt.Println("parse pkg err: ", err)
            return err
        }
    case syscall.EAGAIN:
        // todo 需要把收到的数据缓存拼接一下，那这次recv怎么办，现在不考虑粘包的情况，等调通了来补
    default:
        fmt.Println("read err : ", err)
        return err
    }
    return nil
}
