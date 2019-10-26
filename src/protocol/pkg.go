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
