package pkg

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg          *ini.File
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PageSize     int
)

var ErrorCode = make(map[int]string, 10)
var ErrorMsg = make(map[string]string, 10)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if nil != err {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadErrorCode()
	LoadErrorMsg()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("run_mode").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if nil != err {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("http_port").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("read_timeout").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("write_timeout").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if nil != err {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	PageSize = sec.Key("page_size").MustInt(10)
}

func LoadErrorCode() {
	sec, err := Cfg.GetSection("err_code")
	if nil != err {
		log.Fatalf("Fail to get section 'err_code': %v", err)
	}
	for _, k := range sec.KeyStrings() {
		ki, _ := sec.Key(k).Int()
		ErrorCode[ki] = k
	}
}

func LoadErrorMsg() {
	sec, err := Cfg.GetSection("err_msg")
	if nil != err {
		log.Fatalf("Fail to get section 'err_msg': %v", err)
	}

	for _, k := range sec.KeyStrings() {
		ErrorMsg[k] = sec.Key(k).Value()
	}
}

func GetErrorMsg(code int) string {
	if c, ok := ErrorCode[code]; ok {
		if s, ok := ErrorMsg[c]; ok {
			return s
		}
		return ErrorMsg["error"]
	}
	return ErrorMsg["error"]
}
