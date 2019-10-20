package base

import (
    "fmt"
    "context"
    "syscall"
    "net"
    "sync"
    "time"
    "protocol"
)

type ChannelStatus uint32
const (
    _IdleStatus ChannelStatus = iota
    _RunningStatus
    _StopStatus
)

type TCPChannel struct {
    listener net.Listener
    status ChannelStatus
    quit chan struct{}
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
    channel := &TCPChannel{}
    channel.status = _IdleStatus
    channel.quit = make(chan struct{})
    channel.ip_and_port = addr
    channel.max_buf_len = buf_len
    channel.pending_pkg = make(map[uint32][]byte)
    channel.conn_map = make(map[uint32]net.Conn)
    return channel
}

func (c *TCPChannel) dealConn(ctx context.Context, conn net.Conn, handler func(*protocol.Pkg)error) {
    // 退出的时候要关闭链接
    defer conn.Close()

    c.generate_id++
    new_id := c.generate_id
    c.conn_map[new_id] = conn
    defer delete(c.conn_map, new_id)

    // todo 还没有日志系统，只能这样打了
    fmt.Println("new connect ...from: ", conn.RemoteAddr().String())
    for {
        select {
        case <-c.quit:
            // 收到信号就退出
            return
        default:
            pkg := &protocol.Pkg{}
            err := c.RecvPkg(new_id, pkg)
            if err != nil {
                // todo 还没有日志系统，只能这样打了
                fmt.Println("recv pkg err: ", err)
                // todo 调试需要，先加一句sleep
                time.Sleep(time.Duration(2)*time.Second)
                continue
            }

            // 简单点处理，这里把生成的唯一id给填上
            pkg.Head.Id = new_id
            // 启动一个协程处理，没有context，其实就是没办法结束这个协程
            go handler(pkg)
        }
    }
}

func (c *TCPChannel) Start(wg *sync.WaitGroup, handler func(*protocol.Pkg)error) {

    ctx, cancel := context.WithCancel(context.Background())

    // 如果是协程方式调用，这里可以保证通知到调用的协程这个函数执行完了
    defer func() {
        c.status = _IdleStatus
        cancel()
        wg.Done()
    }()

    // 监听tcp端口
    l, err := net.Listen("tcp", c.ip_and_port)
    if err != nil {
        fmt.Println("listen err : ", err)
        return
    }

    c.listener = l
    c.status = _RunningStatus
    for {
        conn, err := c.listener.Accept()
        if err != nil {
            fmt.Println("accept err : ", err)
            if c.status == _StopStatus {
                fmt.Println("stop accept")
                break
            }
        }
        // 有新链接就启动一个协程处理，这样处理合理么？
        go c.dealConn(ctx, conn, handler)
    }
}

func (c *TCPChannel) ShutDown() {
    if c.status == _RunningStatus {
        c.status = _StopStatus
        close(c.quit)
        c.listener.Close()
    }
}

func (c *TCPChannel) SendPkg(id uint32, pkg *protocol.Pkg) error {
    conn, is_exist := c.conn_map[id]
    if !is_exist {
        return fmt.Errorf("id:%u is not exist", id)
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
    conn, is_exist := c.conn_map[id]
    if !is_exist {
        return fmt.Errorf("id:%u is not exist", id)
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
