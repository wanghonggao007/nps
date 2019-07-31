package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wanghonggao007/nps/lib/common"
	"github.com/wanghonggao007/nps/lib/crypt"
	"github.com/wanghonggao007/nps/lib/daemon"
	"github.com/wanghonggao007/nps/lib/file"
	"github.com/wanghonggao007/nps/lib/install"
	"github.com/wanghonggao007/nps/lib/version"
	"github.com/wanghonggao007/nps/server"
	"github.com/wanghonggao007/nps/server/connection"
	"github.com/wanghonggao007/nps/server/test"
	"github.com/wanghonggao007/nps/server/tool"
	"github.com/wanghonggao007/nps/vender/github.com/astaxie/beego"
	"github.com/wanghonggao007/nps/vender/github.com/astaxie/beego/logs"
	_ "github.com/wanghonggao007/nps/web/routers"

	//
	"strings"
	//
)

var (
	level   string
	logType = flag.String("log", "stdout", "Log output mode（stdout|file）")
)

//
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//
func main() {
	//
	fmt.Println("当前目录:")
	fmt.Println(getCurrentDirectory())
	fmt.Println(common.GetRunPath())
	//
	flag.Parse()
	//beego.LoadAppConfig("ini", filepath.Join(common.GetRunPath(), "conf", "nps.conf"))
	beego.LoadAppConfig("ini", filepath.Join(common.GetRunPath(), "conf", "nps.conf"))
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "test":
			test.TestServerConfig()
			log.Println("test ok, no error")
			return
		case "start", "restart", "stop", "status", "reload":
			daemon.InitDaemon("nps", common.GetRunPath(), common.GetTmpPath())
		case "install":
			install.InstallNps()
			return
		}
	}
	fmt.Println("11111111111111111111111")
	level = beego.AppConfig.String("log_level")
	fmt.Println("22222222222222222222222")
	fmt.Println(level)
	if level = beego.AppConfig.String("log_level"); level == "" {
		level = "7"
	}

	logs.Reset()
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if *logType == "stdout" {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+level+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"level":`+level+`,"filename":"`+beego.AppConfig.String("log_path")+`","daily":false,"maxlines":100000,"color":true}`)
	}
	task := &file.Tunnel{
		Mode: "webServer",
	}
	bridgePort, err := beego.AppConfig.Int("bridge_port")
	if err != nil {
		logs.Error("Getting bridge_port error", err)
		os.Exit(0)
	}
	logs.Info("the version of server is %s ,allow client version to be %s", version.VERSION, version.GetVersion())
	connection.InitConnectionService()
	crypt.InitTls(filepath.Join(common.GetRunPath(), "conf", "server.pem"), filepath.Join(common.GetRunPath(), "conf", "server.key"))
	tool.InitAllowPort()
	tool.StartSystemInfo()
	server.StartNewServer(bridgePort, task, beego.AppConfig.String("bridge_type"))
}
