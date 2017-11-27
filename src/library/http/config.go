package http

import (
    "github.com/BurntSushi/toml"
    "library/file"
    log "github.com/sirupsen/logrus"
    "errors"
)

var (
    ErrorFileNotFound = errors.New("config file not fount")
    ErrorFileParse = errors.New("parse config error")
)

type tcpg struct {
    Mode int     // "1 broadcast" ##(广播)broadcast or  2 (权重)weight
    Name string  // = "group1"
    Filter []string
}
type tcpc struct {
    Listen string
    Port int
}
type TcpConfig struct {
    Enable bool
    Groups map[string]tcpg
    Tcp tcpc
}

type HttpConfig struct {
    Enable bool
    Groups map[string]httpNodeConfig
}

type httpNodeConfig struct {
    Mode int
    Nodes [][]string
    Filter []string
}
type websocketc struct {
    Listen string
    Port int
    Service_ip string
}

const (
    MODEL_BROADCAST = 1  // 广播
    MODEL_WEIGHT    = 2  // 权重
)

const (
    CMD_SET_PRO = 1 // 注册客户端操作，加入到指定分组
    CMD_AUTH    = 2 // 认证（暂未使用）
    CMD_OK      = 3 // 正常响应
    CMD_ERROR   = 4 // 错误响应
    CMD_TICK    = 5 // 心跳包
    CMD_EVENT   = 6 // 事件
    CMD_CONNECT = 7
    CMD_RELOGIN = 8
)

const (
    TCP_MAX_SEND_QUEUE            = 128 //100万缓冲区
    TCP_DEFAULT_CLIENT_SIZE       = 64
    TCP_DEFAULT_READ_BUFFER_SIZE  = 1024
    TCP_RECV_DEFAULT_SIZE         = 4096
    TCP_DEFAULT_WRITE_BUFFER_SIZE = 4096
)

func getHttpConfig() (*tcpc, error) {
    var tcp_config tcpc
    config_file := file.GetCurrentPath()+"/config/admin.toml"
    wfile := file.WFile{config_file}
    if !wfile.Exists() {
        log.Printf("config file %s does not exists", config_file)
        return nil, ErrorFileNotFound
    }

    if _, err := toml.DecodeFile(config_file, &tcp_config); err != nil {
        log.Println(err)
        return nil, ErrorFileParse
    }
    return &tcp_config, nil
}


func getWebsocketConfig() (*websocketc, error) {
    var tcp_config websocketc
    config_file := file.GetCurrentPath()+"/config/admin.toml"
    wfile := file.WFile{config_file}
    if !wfile.Exists() {
        log.Printf("config file %s does not exists", config_file)
        return nil, ErrorFileNotFound
    }

    if _, err := toml.DecodeFile(config_file, &tcp_config); err != nil {
        log.Println(err)
        return nil, ErrorFileParse
    }
    return &tcp_config, nil
}
