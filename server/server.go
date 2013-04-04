package main

import (
	"../util"
	"flag"
	"log"
	"net"
)

var (
	port = flag.String("p", "8010", "Access the service on this port.")
)

func main() {
	flag.Usage = util.Usage
	flag.Parse()
	//flag.Usage()

	server, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", *port))
	if err != nil {
		log.Fatal("CAN'T LISTEN: ", err)
	}
	log.Println("Starting golocproxy on port:", *port)
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal("CAN'T ACCEPT: ", err)
		}
		go newConnect(conn)
	}
}

func newConnect(conn net.Conn) {
	log.Println("Connect:", util.Conn2Str(conn))
	defer conn.Close()

	var buf [1024]byte
	first := true
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Println("Can't Read: ", err)
			break
		}
		if first {
			first = false
			if n == util.TOKEN_LEN {
				token := string(buf[0:n])
				log.Println("token=", token)
				if token == util.C2P_CONNECT {
					//内网服务器启动时连接代理，建立长连接
				} else if token == util.C2P_SESSION {
					//为客户端的单次连接请求建立一个临时的"内网服务器<->代理"的连接
				} else {
					//普通的客户端到代理服务器的连接
				}
			}
		}
		println(string(buf[0:n]))
		conn.Write(buf[0:n])

	}
}
