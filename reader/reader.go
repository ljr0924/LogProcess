package reader

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "os"
    "time"
)

type Reader interface {
    Read(context.Context, chan string)
}

type ReadFromFile struct {
    FilePath string
}

func (r *ReadFromFile) Read(ctx context.Context, ch chan string)  {
    // 打开文件
    f, err := os.Open(r.FilePath)
    if err != nil {
        panic(fmt.Sprintf("open file error, err: %s", err.Error()))
    }
    _, err = f.Seek(0, 2)
    if err != nil {
        panic(fmt.Sprintf("file seek error, err: %s", err.Error()))
    }
    reader := bufio.NewReader(f)
    for {
        select {
        case <- ctx.Done():
            fmt.Println("reader done !!!!!!")
            return
        default:
            // 原始日志读取逻辑
            line, err := reader.ReadBytes('\n')
            if err == io.EOF {
                time.Sleep(time.Second)
                continue
            }

            // 放入读取通道
            ch <- string(line[:len(line)-1])
        }

    }

}