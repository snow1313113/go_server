package base

import (
    "fmt"
    "io"
    "context"
    "syscall"
    "net"
    "sync"
    "time"
    "bytes"
    "encoding/binary"
    "protocol"
)

type ChannelStatus uint32
const (
    _IdleStatus ChannelStatus = iota
    _RunningStatus
    _StopStatus
)

type connectInfo struct {
    // 链接分配的session id
    Id uint32
    // 当前链接
    Connect net.Conn
    // 链接使用的收包buf
    Buffer []byte
}

type TCPChannel struct {
    listener net.Listener
    status ChannelStatus
    quit chan struct{}
    ip_and_port string
    max_buf_len uint32
    // todo 唯一ID怎么定义，简单点的就递增序列分配吧
    generate_id uint32
    // 一个唯一id对应一个connect
    conn_map map[uint32]*connectInfo
}

func NewTCPChannel(addr string, buf_len uint32) *TCPChannel {
    channel := &TCPChannel{}
    channel.status = _IdleStatus
    channel.quit = make(chan struct{})
    channel.ip_and_port = addr
    channel.max_buf_len = buf_len
    channel.conn_map = make(map[uint32]*connectInfo)
    return channel
}

func (c *TCPChannel) dealConn(ctx context.Context, conn net.Conn, handler func(*protocol.Pkg)error) {
    // 退出的时候要关闭链接
    defer conn.Close()

    c.generate_id++
    new_id := c.generate_id
    c.conn_map[new_id] = &connectInfo{Id:new_id, Connect:conn, Buffer:make([]byte, c.max_buf_len)}
    defer delete(c.conn_map, new_id)

    fmt.Println("new connect ...from: ", conn.RemoteAddr().String())
    for {
        select {
        case <-ctx.Done():
            return
        case <-c.quit:
            // 收到信号就退出
            return
        default:
            pkg := &protocol.Pkg{}
            err := c.RecvPkg(new_id, pkg)
            if err != nil {
                if err == io.EOF {
                    // 连接关闭，直接返回
                    fmt.Println("client close connect")
                    return
                }
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
    conn_info, is_exist := c.conn_map[id]
    if !is_exist {
        return fmt.Errorf("id:%u is not exist", id)
    }

    byte_buf, err := pkg.Bytes()
    if err != nil {
        fmt.Println("send err: ", err)
        return err
    }

    _, err = conn_info.Connect.Write(byte_buf)
    if err != nil {
        fmt.Println("write err : ", err)
        return err
    }
    return nil
}

func (c *TCPChannel) RecvPkg(id uint32, pkg *protocol.Pkg) error {
    conn_info, is_exist := c.conn_map[id]
    if !is_exist {
        return fmt.Errorf("id:%u is not exist", id)
    }

    // 先收一个头部
    head_buf := conn_info.Buffer[:protocol.HeadSize]
    _, err := io.ReadFull(conn_info.Connect, head_buf)
    if err != nil {
        return err
    }

    byte_buf := bytes.NewBuffer(head_buf)
    err = binary.Read(byte_buf, binary.LittleEndian, &pkg.Head)
    if err != nil {
        return err
    }

    // 再收body
    body_buf := conn_info.Buffer[protocol.HeadSize:]
    recv_len, err := io.ReadAtLeast(conn_info.Connect, body_buf, int(pkg.Head.BodyLen))
    if err != nil {
        return err
    }
    switch err {
    case nil:
        // todo 这里parse会再解析一遍头部，不优雅，其实还是应该收一个大buffer，然后丢给parse
        _, err = pkg.Parse(conn_info.Buffer[:recv_len + int(protocol.HeadSize)])
        if err != nil {
            return err
        }
    case syscall.EAGAIN:
        // todo 需要把收到的数据缓存拼接一下，那这次recv怎么办，现在不考虑粘包的情况，等调通了来补
    default:
        return err
    }
    return nil
}
