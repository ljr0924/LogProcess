package monitor

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "LogProcess/context"
)

type LogProcessInfo struct {
    ProcessLogNum int32 `json:"process_log_num"`
    ErrorLogNum   int32 `json:"error_log_num"`
    ReadChLen     int   `json:"read_ch_len"`
    WriteChLen    int   `json:"write_ch_len"`
    RunTime       int64 `json:"run_time"`
}

type Monitor struct {
    data      *LogProcessInfo
    StartTime int64 `json:"run_time"`
}

func NewMonitor() *Monitor {
    return &Monitor{
        data:      &LogProcessInfo{},
        StartTime: time.Now().Unix(),
    }
}

func (m *Monitor) StartMonitor(l *context.LogProcessorContext) {
    handler := func(w http.ResponseWriter, r *http.Request) {
        m.data.RunTime = time.Now().Unix() - m.StartTime
        m.data.ErrorLogNum = l.ErrorNum
        m.data.ProcessLogNum = l.ProcessNum
        m.data.ReadChLen = l.GetReadChLen()
        m.data.WriteChLen = l.GetWriteChLen()
        data, err := json.Marshal(m.data)
        if err != nil {
            fmt.Println("json serialize err, ", err.Error())
        }
        _, _ = io.WriteString(w, string(data))
    }
    http.HandleFunc("/monitor", handler)
    _ = http.ListenAndServe(":8000", nil)
}
