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
	proxy.Write([]byte(util.C2P_CONNECT))

	var buf [util.TOKEN_LEN]byte
	for {
		proxy.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := proxy.Read(buf[0:])
	//	proxy.SetReadDeadline(time.Time{})
		if err == nil {
			token := string(buf[0:n])
			if token == util.P2C_NEW_SESSION {
				go session()
			}
		} else {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout(){
				//log.Println("Timeout")
				proxy.SetWriteDeadline(time.Now().Add(2 * time.Second))
				_,werr := proxy.Write([]byte(util.C2P_KEEP_ALIVE)) //send KeepAlive msg
				if(werr != nil){
					log.Println("CAN'T WRITE, err:",werr)
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
	rp.Write([]byte(util.C2P_SESSION))
	lp, err := net.Dial("tcp", *local)
	if err != nil {
		log.Println("Can't' connect:", *local, " err:", err)
		rp.Close()
		return
	}
	go util.CopyFromTo(rp, lp, nil)
	go util.CopyFromTo(lp, rp, nil)
}
