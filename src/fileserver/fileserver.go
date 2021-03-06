package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"net"

	"github.com/astaxie/beego/httplib"
)

func getInternal() {
	// addrs, err := net.InterfaceAddrs()
	// if err != nil {
	// 	os.Stderr.WriteString("Oops:" + err.Error())
	// 	os.Exit(1)
	// }
	// for _, a := range addrs {
	// 	if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 		if ipnet.IP.To4() != nil {
	// 			os.Stdout.WriteString("\t" + ipnet.IP.String() + "\n")
	// 		}
	// 	}
	// }

	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for i, inter := range interfaces {
		fmt.Printf("\t(%d)%s :\n", i, inter.Name)
		addrs, _ := inter.Addrs()
		for indx, addr := range addrs {
			fmt.Printf("\tAddr[%d] = %s\n", indx, addr)
		}
	}
}

func main() {
	port := ":8085"

	arg_num := len(os.Args)
	if arg_num > 1 {
		port = ":" + os.Args[1]
	}

	// 开启web服务
	fmt.Println("/******************局域网内网共享工具********************")
	fmt.Println("\t【文件服务器端口】：")
	fmt.Println("\t", port[1:])
	fmt.Println("\t【本机内网地址】：")
	getInternal()
	fmt.Println("\t【本机公网地址】：")

	//以下代码需要访问公网
	req := httplib.Get("http://getip.stackbang.com")
	req.SetTimeout(time.Second, time.Second)
	externalip, err := req.String()
	if err != nil {
		fmt.Println("\t获取公网地址失败：")
		fmt.Println("\t", err)
	} else {
		if strings.Contains(externalip, "X-Real-IP") {
			strs := strings.Split(externalip, ",")
			fmt.Println("\t" + strs[2][14:])
		} else {
			fmt.Println("\t获取公网地址失败：")
		}
	}

	fmt.Printf("\t在浏览器中访问http://对应网卡的ip%s即可\n", port)
	fmt.Println("********************************************************/")

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	server := http.Server{
		Addr:           port,            // 监听的地址和端口
		Handler:        nil,             // 所有请求需要调用的Handler（实际上这里说是ServeMux更确切）如果为空则设置为DefaultServeMux
		ReadTimeout:    0 * time.Second, // 读的最大Timeout时间
		WriteTimeout:   0 * time.Second, // 写的最大Timeout时间
		MaxHeaderBytes: 256,             // 请求头的最大长度
		TLSConfig:      nil,             // 配置TLS
	}

	errListen := server.ListenAndServe()

	if errListen != nil {
		fmt.Println("ListenAndServer: ", errListen)
	}
	fmt.Println("exit...")
}
