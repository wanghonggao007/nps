package test

import (
	"github.com/wanghonggao007/nps/lib/common"
	"github.com/wanghonggao007/nps/lib/file"
	"github.com/wanghonggao007/nps/vender/github.com/astaxie/beego"
	"log"
	"path/filepath"
	"strconv"
)

func TestServerConfig() {
	var postTcpArr []int
	var postUdpArr []int
	file.GetDb().JsonDb.Tasks.Range(func(key, value interface{}) bool {
		v := value.(*file.Tunnel)
		if v.Mode == "udp" {
			isInArr(&postUdpArr, v.Port, v.Remark, "udp")
		} else if v.Port != 0 {

			isInArr(&postTcpArr, v.Port, v.Remark, "tcp")
		}
		return true
	})
	p, err := beego.AppConfig.Int("web_port")
	if err != nil {
		log.Fatalln("Getting web management port error :", err)
	} else {
		isInArr(&postTcpArr, p, "Web Management port", "tcp")
	}

	if p := beego.AppConfig.String("bridge_port"); p != "" {
		if port, err := strconv.Atoi(p); err != nil {
			log.Fatalln("get Server and client communication portserror:", err)
		} else if beego.AppConfig.String("bridge_type") == "kcp" {
			isInArr(&postUdpArr, port, "Server and client communication ports", "udp")
		} else {
			isInArr(&postTcpArr, port, "Server and client communication ports", "tcp")
		}
	}

	if p := beego.AppConfig.String("httpProxyPort"); p != "" {
		if port, err := strconv.Atoi(p); err != nil {
			log.Fatalln("get http port error:", err)
		} else {
			isInArr(&postTcpArr, port, "https port", "tcp")
		}
	}
	if p := beego.AppConfig.String("https_proxy_port"); p != "" {
		if b, err := beego.AppConfig.Bool("https_just_proxy"); !(err == nil && b) {
			if port, err := strconv.Atoi(p); err != nil {
				log.Fatalln("get https port error", err)
			} else {
				if !common.FileExists(filepath.Join(common.GetRunPath(), beego.AppConfig.String("pemPath"))) {
					log.Fatalf("ssl certFile %s is not exist", beego.AppConfig.String("pemPath"))
				}
				if !common.FileExists(filepath.Join(common.GetRunPath(), beego.AppConfig.String("ketPath"))) {
					log.Fatalf("ssl keyFile %s is not exist", beego.AppConfig.String("pemPath"))
				}
				isInArr(&postTcpArr, port, "http port", "tcp")
			}
		}
	}
}

func isInArr(arr *[]int, val int, remark string, tp string) {
	for _, v := range *arr {
		if v == val {
			log.Fatalf("the port %d is reused,remark: %s", val, remark)
		}
	}
	if tp == "tcp" {
		if !common.TestTcpPort(val) {
			log.Fatalf("open the %d port error ,remark: %s", val, remark)
		}
	} else {
		if !common.TestUdpPort(val) {
			log.Fatalf("open the %d port error ,remark: %s", val, remark)
		}
	}
	*arr = append(*arr, val)
	return
}
