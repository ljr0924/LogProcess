package context

import (
    "LogProcess/processor"
    "LogProcess/reader"
    "LogProcess/writer"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/signal"
    "time"
)

type LogProcessorContext struct {
    readCh    chan string
    writeCh   chan []string
    ctx       context.Context
    cancel    context.CancelFunc
    reader    reader.Reader
    writer    writer.Writer
    processor processor.Processor
    monitor   *Monitor
}

type Monitor struct {
    ProcessLogNum    int   `json:"process_log_num"`
    WriteLogNum      int   `json:"write_log_num"`
    ReadErrLogNum    int   `json:"read_err_log_num"`
    ProcessErrLogNum int   `json:"process_err_log_num"`
    WriteErrLogNum   int   `json:"write_err_log_num"`
    LastTime         int64 `json:"last_time"`
}

func (m *Monitor) AddWriteErrLogNum() {
    m.WriteErrLogNum++
    m.LastTime = time.Now().Unix()
}

func (m *Monitor) AddProcessErrLogNum() {
    m.ProcessErrLogNum++
    m.LastTime = time.Now().Unix()
}

func (m *Monitor) AddReadErrLogNum() {
    m.ReadErrLogNum++
    m.LastTime = time.Now().Unix()
}

func (m *Monitor) AddProcessLogNum() {
    m.ProcessLogNum++
    m.LastTime = time.Now().Unix()
}

func (m *Monitor) AddWriteLogNum() {
    m.WriteLogNum++
    m.LastTime = time.Now().Unix()
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
        monitor:   &Monitor{},
    }
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
                l.monitor.AddWriteErrLogNum()
                continue
            }
            l.monitor.AddWriteLogNum()
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
                l.monitor.AddProcessErrLogNum()
                continue
            }
            l.writeCh <- logDataArr
            // 监控
            l.monitor.AddProcessLogNum()
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

func (l *LogProcessorContext) Monitor() {
    handler := func(w http.ResponseWriter, r *http.Request) {
        data, err := json.Marshal(l.monitor)
        if err != nil {
            fmt.Println("json serialize err, ", err.Error())
        }
        _, _ = io.WriteString(w, string(data))
    }
    http.HandleFunc("/monitor", handler)
    _ = http.ListenAndServe(":8000", nil)
}
