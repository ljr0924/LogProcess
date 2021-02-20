package context

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "sync/atomic"

    "LogProcess/processor"
    "LogProcess/reader"
    "LogProcess/writer"
)

type LogProcessorContext struct {
    readCh     chan string
    writeCh    chan []string
    ProcessNum int32
    ErrorNum   int32
    ctx        context.Context
    cancel     context.CancelFunc
    reader     reader.Reader
    writer     writer.Writer
    processor  processor.Processor
}

func (l *LogProcessorContext) AddErrorNum() {
    atomic.AddInt32(&l.ErrorNum, 1)
}

func (l *LogProcessorContext) AddProcessNum() {
    atomic.AddInt32(&l.ProcessNum, 1)
}

func NewLogContext(r reader.Reader, w writer.Writer, p processor.Processor) *LogProcessorContext {

    ctx, cancel := context.WithCancel(context.TODO())

    return &LogProcessorContext{
        readCh:    make(chan string),
        writeCh:   make(chan []string),
        ctx:       ctx,
        cancel:    cancel,
        reader:    r,
        writer:    w,
        processor: p,
    }
}

func (l *LogProcessorContext) GetReadChLen() int {
    return len(l.readCh)
}

func (l *LogProcessorContext) GetWriteChLen() int {
    return len(l.writeCh)
}

func (l *LogProcessorContext) Read() {
    l.reader.Read(l.ctx, l.readCh)
}

func (l *LogProcessorContext) Write() {
    for {
        select {
        case <-l.ctx.Done():
            fmt.Println("writer done !!!!!!!!")
            return
        case m := <-l.writeCh:
            err := l.writer.Write(m)
            if err != nil {
                l.AddErrorNum()
                continue
            }
            l.AddProcessNum()
        }
    }

}

func (l *LogProcessorContext) Process() {
    for {
        select {
        case <-l.ctx.Done():
            fmt.Println("processor done!!!!!!")
            return
        case m := <-l.readCh:
            // 日志处理逻辑
            logDataArr, err := l.processor.Process(m)
            if err != nil {
                l.AddErrorNum()
                continue
            }
            l.writeCh <- logDataArr
            // 监控
            l.AddProcessNum()
        }
    }
}

func (l *LogProcessorContext) Shutdown() {
    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt)
    switch <-done {
    case os.Interrupt:
        fmt.Println("ready to shutdown")
        l.cancel()
    }
}
