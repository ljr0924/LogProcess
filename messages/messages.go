package messages

type NginxLog struct {
    Time         string `bson:"time"`
    Host         string `bson:"host"`
    SourceIp     string `bson:"source_ip"`
    Method       string `bson:"method"`
    Uri          string `bson:"uri"`
    Protocol     string `bson:"protocol"`
    ResponseTime string `bson:"response_time"`
    Bytes        string `bson:"bytes"`
    Code         string `bson:"code"`
    Dip          string `bson:"dip"`
    UpstreamTime string `bson:"upstream_time"`
    Agent        string `bson:"agent"`
}
