package main

import (
	//"database/sql"
	"library"
	"library/services"
	//"github.com/siddontang/go-mysql/canal"

	//"library/base"
	//"library/workers"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	//"strconv"
	//"subscribe"
	"os"
	"os/signal"
	"syscall"
	_"net/http/pprof"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
)


func writePid() {
	var data_str = []byte(fmt.Sprintf("%d", os.Getpid()));
	ioutil.WriteFile(library.GetCurrentPath() + "/wing-binlog-go.pid", data_str, 0777)  //写入文件(字节数组)
}

func killPid() {
	dat, _ := ioutil.ReadFile(library.GetCurrentPath() + "/wing-binlog-go.pid")
	fmt.Print(string(dat))
	pid, _ := strconv.Atoi(string(dat))
	log.Println("给进程发送终止信号：", pid)
	err := syscall.Kill(pid, syscall.SIGTERM)
	log.Println(err)
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "stop" {
		killPid()
		return
	}

	writePid()

	//标准输出重定向
	//library.Reset()
	go func() {
		//http://localhost:6060/debug/pprof/  内存性能分析工具
		//go tool pprof logDemo.exe --text a.prof
		//go tool pprof your-executable-name profile-filename
		//go tool pprof your-executable-name http://localhost:6060/debug/pprof/heap
		//go tool pprof wing-binlog-go http://localhost:6060/debug/pprof/heap
		//https://lrita.github.io/2017/05/26/golang-memory-pprof/
		//然后执行 text
		//go tool pprof -alloc_space http://127.0.0.1:6060/debug/pprof/heap
		//top20 -cum
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wing_log := library.GetLogInstance()
	//释放日志资源
	defer library.FreeLogInstance()

	/*file := &library.WFile{"C:\\__test.txt"}
	str := file.ReadAll()
	//if str != "123" {
		log.Println("ReadAll error: ==>" + str + "<==", len(str))
	//}
	return*/

	log.SetFlags(log.LstdFlags | log.Lshortfile)


	cpu := runtime.NumCPU()
	wing_log.Println("cpu num: ", cpu)

	//指定cpu为多核运行
	runtime.GOMAXPROCS(cpu)

	current_path := library.GetCurrentPath()
	wing_log.Println(current_path)

	config_file := current_path + "/config/mysql.toml"
	tcp_config_file := current_path + "/config/tcp.toml"

	//config_obj := &library.Ini{config_file}
	//config := config_obj.Parse()
	//if config == nil {
	//	wing_log.Println("read config file: " + config_file + " error")
	//	return
	//}
	//wing_log.Println(config)


	//config_file := "/tmp/__test_mysql.toml"
	config := &library.WConfig{config_file}

	app_config, err:= config.GetMysql()

	if err != nil {
		log.Println(err)
		return
	}

	wtcp_config := &library.WConfig{tcp_config_file}
	//{map[1:{1 group1} 2:{2 group2}] {0.0.0.0 9998}}
	tcp_config, err := wtcp_config.GetTcp()
	if err != nil {
		log.Println(err)
		return
	}
	//user := string(config["mysql"]["user"].(string))
	//password := string(config["mysql"]["password"].(string))
	//port := string(config["mysql"]["port"].(string))
	//host := string(config["mysql"]["host"].(string))

	//slave_id_str := string(config["client"]["slave_id"].(string))
	//slave_id, _ := strconv.Atoi(slave_id_str)

	//db_name := string(config["mysql"]["db_name"].(string))
	//charset := string(config["mysql"]["charset"].(string))
	//db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+db_name+"?charset="+charset)

	//if nil != err {
	//	wing_log.Println(err)
	//	return
	//}

	//defer db.Close()

	tcp_service := services.NewTcpService(tcp_config.Tcp.Listen, tcp_config.Tcp.Port, tcp_config)
	websocket_service := services.NewWebSocketService("0.0.0.0", 9997)

	tcp_service.Start()

	blog := library.Binlog{DB_Config:app_config}
	defer blog.Close()
	blog.Start(tcp_service, websocket_service)

	<-sc

	//redis := &subscribe.Redis{}
	//tcp := &subscribe.Tcp{}
    //
	////subscribes
	//notify := []base.Subscribe{redis, tcp}
	//binlog := &workers.Binlog{}
    //
	//defer binlog.End(notify)
    //
	//binlog.Start(notify)
	//binlog.Loop(notify)
}
