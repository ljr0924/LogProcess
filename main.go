package main

import (
    "LogProcess/context"
    "LogProcess/processor"
    "LogProcess/reader"
    "LogProcess/writer"
    "fmt"
    "os"
    "time"
)

func main() {
    r := &reader.ReadFromFile{"./access.log"}
    w := &writer.WriteToMongo{writer.CollNginxLog}
    p := &processor.NginxProcessor{}
    logCtx := context.NewLogContext(r, w, p)

    go logCtx.Read()
    go logCtx.Process()
    go logCtx.Write()
    go logCtx.Monitor()

    // 模拟日志输入
    go mockNginxLog(100)

    logCtx.Shutdown()

}

func mockNginxLog(logNum int) {
    time.Sleep(time.Second)
    f, err := os.OpenFile("./access.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {
        fmt.Println("write err: ", err.Error())
        return
    }
    count := 0
    for count < logNum {
        var log string
        if count % 2 == 0 {
            // 错误的日志
            log = "[18/Feb/2017:19:16:59] | test.xyz | 115.33.60.172 | 12313 | POST | /api/1.1/device/info | HTTP/1.1 | 0.003 | 43 | 200 | 127.0.0.1:6000 | 0.003 | 200 | \"-\" | \"Apache-HttpClient/UNAVAILABLE(java 1.4)\"\n"
        } else {
            // 正常的日志
            log = "[18/Feb/2017:19:16:59] | test.xyz | 115.33.60.172 | POST | /api/1.1/device/info | HTTP/1.1 | 0.003 | 43 | 200 | 127.0.0.1:6000 | 0.003 | 200 | \"-\" | \"Apache-HttpClient/UNAVAILABLE(java 1.4)\"\n"
        }
        _, err = f.Write([]byte(log))
        if err != nil {
            fmt.Println("write err: ", err.Error())
        }
        time.Sleep(time.Second)
        count++
    }
    _ = f.Close()


}
