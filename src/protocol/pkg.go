package protocol

import (
    "bytes"
    "errors"
    "encoding/binary"
)

var (
    ErrParseHeadInvalid = errors.New("parse pkg head len is invalid")
    ErrParseBodyBufferNotEnough = errors.New("parse pkg body len is not enough")
    ErrParseDataBodyInvalid = errors.New("parse pkg body data len is invalid")
    ErrByteDataBodyInvalid = errors.New("bytes pkg body data len is invalid")
)

var (
    HeadSize = uint32(4 + 4 + 4 + 4 + 4)
)

// 需要保证头部是定长整数
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

func (pkg *Pkg) Parse(data []byte) ([]byte, error) {
    if len(data) < int(HeadSize) {
        return data, ErrParseHeadInvalid
    }
    buffer := bytes.NewBuffer(data)
    err := binary.Read(buffer, binary.LittleEndian, &pkg.Head)
    if err != nil {
        return data, err
    }
    remain := buffer.Bytes()
    if pkg.Head.BodyLen > uint32(buffer.Len()) {
        return remain, ErrParseBodyBufferNotEnough
    }
    pkg.Body = remain[:pkg.Head.BodyLen]
    return remain[pkg.Head.BodyLen:], nil
}

func (pkg *Pkg) Bytes() ([]byte, error) {
    if pkg.Body != nil {
        if pkg.Head.BodyLen != uint32(len(pkg.Body)) {
            return nil, ErrByteDataBodyInvalid
        }
    } else {
        if pkg.Head.BodyLen != 0 {
            return nil, ErrByteDataBodyInvalid
        }
    }

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
