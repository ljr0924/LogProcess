package writer

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    kMongoUri     = "mongodb://%s:%s"
    kMongoServer  = "127.0.0.1"
    kMongoPort    = "27017"
    kDataBaseTest = "test"
    kCollNginxLog = "nginx_log"
)

var Client *mongo.Client
var CollNginxLog *mongo.Collection

func GetCollNginxLog() *mongo.Collection {

    return Client.Database(kDataBaseTest).Collection(kCollNginxLog)

}

func connect() error {

    ctx := context.TODO()

    clientOptions := options.Client().ApplyURI(fmt.Sprintf(kMongoUri, kMongoServer, kMongoPort))
    var err error
    Client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        fmt.Println(err)
        return err
    }

    err = Client.Ping(context.TODO(), nil)
    if err != nil {
        fmt.Println(err)
        return err
    }

    return nil
}

func init() {
    fmt.Println("正在初始化mongo连接对象")
    err := connect()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("mongo连接对象初始化成功")
    CollNginxLog = GetCollNginxLog()
}
