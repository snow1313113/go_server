package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "base"
    "utils"
    "protocol"
)

var (
    log *utils.Logger
)

func main() {
    fmt.Println("start")

    sig_chan := make(chan os.Signal, 1)
    signal.Notify(sig_chan, syscall.SIGINT, syscall.SIGTERM)

    utils.Daemon(1, 1)

    log, _ = utils.NewLogger(utils.DebugLevel, "")
	defer log.Close()

    log.Debug("new server")

    svr := base.NewServer("localhost:1234", 512)
    // 注意了，ExampleService的方法是用指针作为reciver的，所以这里要传指针给空接口，
    // 不然反射出来的就不是ExampleService了
    err := svr.Registered(&ExampleService{}, protocol.ExampleService_Desc)
    if err != nil {
        log.Error("regist service err: %v", err)
        return
    }
    log.Info("registred succeed")

    go svr.Run()

    // 等待收到退出信号
    for sig := range sig_chan {
        log.Info("recv signal: %v", sig)
        break
    }

    // 直接调用stop
    svr.Stop()
    log.Info("svr stop")
}

