package writer

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"

    "LogProcess/messages"
)

type Writer interface {
    Write([]string) error
}

type WriteToMongo struct {
    Coll *mongo.Collection
}

func (w *WriteToMongo) Write(logData []string) error {
    // 日志数据库写入逻辑
    nginxLog := messages.NginxLog{
        Time:         logData[0],
        Host:         logData[1],
        SourceIp:     logData[2],
        Method:       logData[3],
        Uri:          logData[4],
        Protocol:     logData[5],
        ResponseTime: logData[6],
        Bytes:        logData[7],
        Code:         logData[8],
        Dip:          logData[9],
        UpstreamTime: logData[10],
        Agent:        logData[13],
    }
    _, err := w.Coll.InsertOne(context.TODO(), nginxLog)
    if err != nil {
        return err
    }
    return nil
}
