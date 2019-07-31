package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wanghonggao007/nps/client"
	"github.com/wanghonggao007/nps/lib/common"
	"github.com/wanghonggao007/nps/lib/config"
	"github.com/wanghonggao007/nps/lib/daemon"
	"github.com/wanghonggao007/nps/lib/file"
	"github.com/wanghonggao007/nps/lib/version"
	"github.com/wanghonggao007/nps/vender/github.com/astaxie/beego/logs"
	"github.com/wanghonggao007/nps/vender/github.com/ccding/go-stun/stun"
)

var (
	serverAddr   = flag.String("server", "", "Server addr (ip:port)")
	configPath   = flag.String("config", "", "Configuration file path")
	verifyKey    = flag.String("vkey", "", "Authentication key")
	logType      = flag.String("log", "stdout", "Log output mode（stdout|file）")
	connType     = flag.String("type", "tcp", "Connection type with the server（kcp|tcp）")
	proxyUrl     = flag.String("proxy", "", "proxy socks5 url(eg:socks5://111:222@127.0.0.1:9007)")
	logLevel     = flag.String("log_level", "7", "log level 0~7")
	registerTime = flag.Int("time", 2, "register time long /h")
	localPort    = flag.Int("local_port", 2000, "p2p local port")
	password     = flag.String("password", "", "p2p password flag")
	target       = flag.String("target", "", "p2p target")
	localType    = flag.String("local_type", "p2p", "p2p target")
	logPath      = flag.String("log_path", "npc.log", "npc log path")
)

func main() {
	os.Args = []string{"w"}
	flag.Parse()

	if len(os.Args) >= 2 {
		fmt.Println("参数大于2", len(os.Args), os.Args[1])
		switch os.Args[1] {
		case "status":
			if len(os.Args) > 2 {
				path := strings.Replace(os.Args[2], "-config=", "", -1)
				fmt.Println("替换参数：", path)
				client.GetTaskStatus(path)
			}
		case "register":
			flag.CommandLine.Parse(os.Args[2:])
			fmt.Println("命令行参数：", os.Args[0:])
			fmt.Println("服务器参数：", *serverAddr)
			fmt.Println("服务器参数：", *verifyKey)
			fmt.Println("服务器参数：", *connType)
			fmt.Println("服务器参数：", *proxyUrl)
			fmt.Println("服务器参数：", *registerTime)
			client.RegisterLocalIp(*serverAddr, *verifyKey, *connType, *proxyUrl, *registerTime)
		case "nat":
			fmt.Println("nat")
			nat, host, err := stun.NewClient().Discover()
			fmt.Println(nat)
			fmt.Println(host)
			fmt.Println(err)
			if err != nil || host == nil {
				logs.Error("get nat type error", err)
				return
			}
			fmt.Printf("nat type: %s \npublic address: %s\n", nat.String(), host.String())
			os.Exit(0)
		}
	}
	fmt.Println("参数小于2", len(os.Args))
	daemon.InitDaemon("npc", common.GetRunPath(), common.GetTmpPath())
	fmt.Println("路径", common.GetRunPath(), common.GetTmpPath())
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if *logType == "stdout" {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+*logLevel+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"level":`+*logLevel+`,"filename":"`+*logPath+`","daily":false,"maxlines":100000,"color":true}`)
	}
	//p2p or secret command
	if *password != "" {
		commonConfig := new(config.CommonConfig)
		commonConfig.Server = *serverAddr
		commonConfig.VKey = *verifyKey
		commonConfig.Tp = *connType
		localServer := new(config.LocalServer)
		localServer.Type = *localType
		localServer.Password = *password
		localServer.Target = *target
		localServer.Port = *localPort
		commonConfig.Client = new(file.Client)
		commonConfig.Client.Cnf = new(file.Config)
		client.StartLocalServer(localServer, commonConfig)
		return
	}
	env := common.GetEnvMap()
	if *serverAddr == "" {
		*serverAddr, _ = env["NPC_SERVER_ADDR"]
	}
	if *verifyKey == "" {
		*verifyKey, _ = env["NPC_SERVER_VKEY"]
	}
	logs.Info("the version of client is %s, the core version of client is %s", version.VERSION, version.GetVersion())
	if *verifyKey != "" && *serverAddr != "" && *configPath == "" {
		for {
			client.NewRPClient(*serverAddr, *verifyKey, *connType, *proxyUrl, nil).Start()
			logs.Info("It will be reconnected in five seconds")
			time.Sleep(time.Second * 5)
		}
	} else {
		if *configPath == "" {
			//*configPath = "npc.conf"
			*configPath = "conf\\npc.conf"
		}
		client.StartFromFile(*configPath)
	}
}
