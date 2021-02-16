package processor

import (
    "errors"
    "strings"
)

type Processor interface {
    Process(string) ([]string, error)
}

type NginxProcessor struct {}

func (r *NginxProcessor) Process(log string) ([]string, error) {
    logDataArr := strings.Split(log, " | ")
    if len(logDataArr) != 14 {
        return nil, errors.New("log process err")
    }
    return logDataArr, nil
}
