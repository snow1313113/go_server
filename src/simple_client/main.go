package main

import (
    "fmt"
    "net"
    "io"
    "syscall"
    "time"
    "context"
    "sync"
    "bytes"
    "encoding/binary"
    "protocol"
    pb "github.com/golang/protobuf/proto"
)

func sendPkg(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
    defer wg.Done()

    for i := 1; i < 10; i++ {
        select {
        case <-ctx.Done():
            return
        default:
            pkg := protocol.Pkg{}

            pkg.Head.Id = 0
            pkg.Head.Cmd = 0x01
            pkg.Head.Seq = uint32(i)

            req := protocol.HelloReq {"test hello", uint32(i*100)}
            body, err := pb.Marshal(&req)
            if err != nil {
                fmt.Println("Marshal err:", err)
                return
            }
            pkg.Body = body
            pkg.Head.BodyLen = uint32(len(pkg.Body))

            pkg_buf, err := pkg.Bytes()
            if err != nil {
                fmt.Println("pkg bytes err:", err)
                return
            }

            pkg_len := len(pkg_buf)
            buffer := new(bytes.Buffer)
            err = binary.Write(buffer, binary.LittleEndian, uint32(pkg_len))
            if err != nil {
                fmt.Println("binary.write pkg_len err:", err)
                return
            }
            binary.Write(buffer, binary.LittleEndian, pkg_buf)

            _, err = conn.Write(buffer.Bytes())
            if err != nil {
                fmt.Println("write err : ", err)
                return
            }

            time.Sleep(time.Duration(2)*time.Second)
        }
    }
}

func recvPkg(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
    defer wg.Done()

    for {
        select {
        case <-ctx.Done():
            return
        default:
            len_buf := make([]byte, 4)
            _, err := io.ReadFull(conn, len_buf)
            if err != nil {
                fmt.Println("ReadFull err : ", err)
                return
            }

            pkg_len := binary.LittleEndian.Uint32(len_buf)

            buf := make([]byte, pkg_len)
            buf_len, err := io.ReadFull(conn, buf)
            switch err {
            case nil:
                pkg := &protocol.Pkg{}
                err = pkg.Parse(buf[:buf_len])
                if err != nil {
                    fmt.Println("parse err: ", err)
                    return
                }
                fmt.Println("recv rsp: ", pkg)
            case io.EOF:
                // 连接关闭，直接返回
                fmt.Println("client close connect")
                return
            case syscall.EAGAIN:
                // todo 需要把收到的数据缓存拼接一下，那这次recv怎么办，现在不考虑粘包的情况，等调通了来补
            default:
                time.Sleep(time.Duration(1)*time.Second)
                fmt.Println("recv err: ", err)
            }
        }
    }
}

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

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    wg := sync.WaitGroup{}

    wg.Add(1)
    go sendPkg(ctx, &wg, conn)

    wg.Add(1)
    go recvPkg(ctx, &wg, conn)

    wg.Wait()
}
