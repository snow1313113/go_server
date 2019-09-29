package simple_worker

type Worker interface {
    Init() int
    Run() string
    Stop() bool
}

type EchoWorker struct {
    Word string
}

func (w  EchoWorker) Init() int {
    return 0
}

func (w  EchoWorker) Run() string {
    return w.Word
}

func (w  EchoWorker) Stop() bool {
    return true
}

