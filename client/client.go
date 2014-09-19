package main

import (
	"../util"
	"flag"
	//	"io"
	"log"
	"net"
	"time"
)

var (
	local  = flag.String("l", "127.0.0.1:80", "Address of the local app service")
	remote = flag.String("r", "127.0.0.1:8010", "Address of the golocproxy server")
	pwd    = flag.String("pwd", "jjg", "password to access server")
)

func main() {
	flag.Usage = util.Usage
	flag.Parse()

	log.Println("golocproxy client starting: ", *local, "->", *remote)
	for {
		connectServer()
		time.Sleep(10 * time.Second) //retry after 10s
	}
}

func connectServer() {
	proxy, err := net.DialTimeout("tcp", *remote, 5*time.Second)
	if err != nil {
		log.Println("CAN'T CONNECT:", *remote, " err:", err)
		return
	}
	defer proxy.Close()
	util.WriteString(proxy, *pwd+"\n"+util.C2P_CONNECT)

	for {
		proxy.SetReadDeadline(time.Now().Add(2 * time.Second))
		msg, err := util.ReadString(proxy)
		//	proxy.SetReadDeadline(time.Time{})
		if err == nil {
			if msg == util.P2C_NEW_SESSION {
				go session()
			} else {
				log.Println(msg)
			}
		} else {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				//log.Println("Timeout")
				proxy.SetWriteDeadline(time.Now().Add(2 * time.Second))
				_, werr := util.WriteString(proxy, util.C2P_KEEP_ALIVE) //send KeepAlive msg
				if werr != nil {
					log.Println("CAN'T WRITE, err:", werr)
					return
				}

				continue
			} else {
				log.Println("SERVER CLOSE, err:", err)
				return
			}
		}
		//time.Sleep(2*time.Second)
	}

}

//客户端单次连接处理
func session() {
	log.Println("Create Session")
	rp, err := net.Dial("tcp", *remote)
	if err != nil {
		log.Println("Can't' connect:", *remote, " err:", err)
		return
	}
	//defer util.CloseConn(rp)
	util.WriteString(rp, *pwd+"\n"+util.C2P_SESSION)
	lp, err := net.Dial("tcp", *local)
	if err != nil {
		log.Println("Can't' connect:", *local, " err:", err)
		rp.Close()
		return
	}
	go util.CopyFromTo(rp, lp, nil)
	go util.CopyFromTo(lp, rp, nil)
}
