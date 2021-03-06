package binlog

import (
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"errors"
	"library/file"
	"library/services"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"time"
	"unicode/utf8"
	"os"
	"context"
	"sync"
)

var (
	ErrorFileNotFound = errors.New("文件不存在")
	ErrorFileParse = errors.New("配置解析错误")
)

type AppConfig struct {
	Addr     string `toml:"addr"`
	User     string `toml:"user"`
	Password string `toml:"password"`

	Charset         string        `toml:"charset"`
	ServerID        uint32        `toml:"server_id"`
	Flavor          string        `toml:"flavor"`
	HeartbeatPeriod time.Duration `toml:"heartbeat_period"`
	ReadTimeout     time.Duration `toml:"read_timeout"`

	BinFile string `toml:"bin_file"`
	BinPos int64  `toml:"bin_pos"`
}

//type ClientConfig struct {
//	Slave_id int
//	Ignore_tables []string
//	Bin_file string
//	Bin_pos int64
//}

//type MysqlConfig struct {
//	Host string
//	User string
//	Password string
//	Port int
//	Charset string
//	DbName string
//}

type Binlog struct {
	Config *AppConfig
	handler *canal.Canal
	isClosed bool
	BinlogHandler binlogHandler
	ctx *context.Context
	wg *sync.WaitGroup
	lock *sync.Mutex                      // 互斥锁，修改资源时锁定
}

type positionCache struct {
	pos mysql.Position
	index int64
}

const (
	MAX_CHAN_FOR_SAVE_POSITION = 128
	defaultBufSize = 4096
	DEFAULT_FLOAT_PREC = 6
)

type binlogHandler struct {
	Event_index int64
	canal.DummyEventHandler
	buf               []byte
	services map[string] services.Service
	servicesCount int
	cacheHandler *os.File
	lock *sync.Mutex                      // 互斥锁，修改资源时锁定
	wg *sync.WaitGroup
	isClosed bool
	ctx *context.Context
}

// 获取mysql配置
func GetMysqlConfig() (*AppConfig, error) {
	var app_config AppConfig
	config_file := file.GetCurrentPath() + "/config/canal.toml"
	wfile := file.WFile{config_file}
	if !wfile.Exists() {
		log.Errorf("配置文件%s不存在 %s", config_file)
		return nil, ErrorFileNotFound
	}
	if _, err := toml.DecodeFile(config_file, &app_config); err != nil {
		log.Println(err)
		return nil, ErrorFileParse
	}
	return &app_config, nil
}

var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}
var hex = "0123456789abcdef"

