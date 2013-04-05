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
	chSession := make(chan net.Conn)
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal("CAN'T ACCEPT: ", err)
		}
		go onConnect(conn, chSession)
	}
}

func onConnect(conn net.Conn, chSession chan net.Conn) {
	strConn := util.Conn2Str(conn)
	log.Println("Connect:", strConn)
	//defer func() {
	//	log.Println("Close:", strConn)
	//	conn.Close()
	//}()
	var buf [util.TOKEN_LEN]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		log.Println("Can't Read: ", err)
		return
	}
	if n == util.TOKEN_LEN {
		token := string(buf[0:n])
		//log.Println("token=", token)
		if token == util.C2P_CONNECT {
			//内网服务器启动时连接代理，建立长连接
			initC2PConnect(conn)
			return
		} else if token == util.C2P_SESSION {
			//为客户端的单次连接请求建立一个临时的"内网服务器<->代理"的连接
			initClientSession(conn, chSession)
			return
		}
	}
	//普通的客户端到代理服务器的连接
	notifyClientCreateSession(conn, buf[0:n], chSession)
	//println(string(buf[0:n]))
	//conn.Write(buf[0:n])

}

var clientProxy net.Conn = nil

func initC2PConnect(conn net.Conn) {
	defer conn.Close()
	if clientProxy != nil {
		conn.Write([]byte("P2C:service existing"))
		return
	}
	println("REG service")
	clientProxy = conn
	var buf [1]byte
	for {
		_, err := clientProxy.Read(buf[0:])
		if err != nil {
			println("UNREG service")
			clientProxy = nil
			break
		}
	}
}
func initClientSession(conn net.Conn, chSession chan net.Conn) {
	chSession <- conn
}
func notifyClientCreateSession(conn net.Conn, bufReaded []byte, chSession chan net.Conn) {
	if clientProxy == nil {
		conn.Write([]byte("NO SERVICE"))
		conn.Close()
		return
	}
	clientProxy.Write([]byte(util.P2C_NEW_SESSION))
	connSession := <-chSession
	log.Println("Start transfer...")
	connSession.Write(bufReaded)
	go util.CopyFromTo(conn, connSession)
	go util.CopyFromTo(connSession, conn)
}
