package protocol

import (
    "bytes"
    "errors"
    "encoding/binary"
)

// todo 可以自定义个error struct更好一点
var (
    ErrDataBodyInvalid = errors.New("pkg body data len is invalid")
)

type PkgHead struct {
    Id uint32
    Cmd uint32
    Seq uint32
    Ret int32
    BodyLen uint32
}

type Pkg struct {
    Head PkgHead
    Body []byte
}

func (pkg *Pkg) Parse(data []byte) error {
    buffer := bytes.NewBuffer(data)
    err := binary.Read(buffer, binary.LittleEndian, &pkg.Head)
    if err != nil {
        return err
    }
    if pkg.Head.BodyLen != uint32(buffer.Len()) {
        return ErrDataBodyInvalid
    }
    pkg.Body = buffer.Bytes()
    return nil
}

func (pkg *Pkg) Bytes() ([]byte, error) {
    buffer := new(bytes.Buffer)
    err := binary.Write(buffer, binary.LittleEndian, pkg.Head)
    if err != nil {
        return nil, err
    }
    if pkg.Body != nil {
        buffer.Write(pkg.Body)
    }
    return buffer.Bytes(), nil
}

// todo 其实对于普通的struct结构来说，parse和bytes都是一样的，
// 所以其实看看能不能根据反射直接写成两个通用函数了
type Message interface {
    Parse([]byte) error
    Bytes() ([]byte, error)
}

/// 测试用的几个协议包结构
type HelloReq struct {
    Name uint32
    Num uint32
}

func (h *HelloReq) Parse(data []byte) error {
    buffer := bytes.NewBuffer(data)
    err := binary.Read(buffer, binary.LittleEndian, h)
    if err != nil {
        return err
    }
    return nil
}

func (h *HelloReq) Bytes() ([]byte, error) {
    buffer := new(bytes.Buffer)
    err := binary.Write(buffer, binary.LittleEndian, h)
    if err != nil {
        return nil, err
    }
    return buffer.Bytes(), nil
}

type HelloRsp struct {
    Name uint32
    Num uint32
    Seq uint32
}

func (h *HelloRsp) Parse(data []byte) error {
    buffer := bytes.NewBuffer(data)
    err := binary.Read(buffer, binary.LittleEndian, h)
    if err != nil {
        return err
    }
    return nil
}

func (h *HelloRsp) Bytes() ([]byte, error) {
    buffer := new(bytes.Buffer)
    err := binary.Write(buffer, binary.LittleEndian, *h)
    if err != nil {
        return nil, err
    }
    return buffer.Bytes(), nil
}



