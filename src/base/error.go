package base

type RpcError struct {
    Code int32
    Msg string
}

func NewRpcError(code int32, msg string) error {
    return &RpcError{code, msg}
}

func (err *RpcError) Error() string {
    return err.Msg
}
