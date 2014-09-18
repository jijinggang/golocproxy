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
	port = flag.String("p", "8010", "The Listen port of golocproxy, users will access the port.")
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
	chSession := make(chan net.Conn, 100)
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Can't Accept: ", err)
			continue
		}
		go onConnect(conn, chSession)
	}
}

func onConnect(conn net.Conn, chSession chan net.Conn) {
	strConn := util.Conn2Str(conn)
	log.Println("Connect:", strConn)
	var buf [util.TOKEN_LEN]byte
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf[0:])
	conn.SetReadDeadline(time.Time{})
	//println("Read:", string(buf[0:n]))
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			log.Println("Timeout:", string(buf[0:n]), err)
			userConnect(conn, buf[0:n], chSession)
			return
		}
		log.Println("Can't Read: ", err)
		conn.Close()
		return
	}
	if n == util.TOKEN_LEN {
		token := string(buf[0:n])
		//log.Println("token=", token)
		if token == util.C2P_CONNECT {
			//内网服务器启动时连接代理，建立长连接
			clientConnect(conn)
			return
		} else if token == util.C2P_SESSION {
			//为客户端的单次连接请求建立一个临时的"内网服务器<->代理"的连接
			initUserSession(conn, chSession)
			return
		}
	}
	//普通的客户端到代理服务器的连接
	userConnect(conn, buf[0:n], chSession)
	//println(string(buf[0:n]))
	//conn.Write(buf[0:n])

}

//代理客户端连接
var clientProxy net.Conn = nil

//处理golocproxy client的连接
func clientConnect(conn net.Conn) {
	defer util.CloseConn(conn) // conn.Close()
	if clientProxy != nil {
		conn.Write([]byte("SERVICE EXIST"))
		return
	}
	println("REG SERVICE")
	clientProxy = conn
	defer func() {
		clientProxy = nil
	}()
	var buf [util.TOKEN_LEN]byte
	for {
		_, err := clientProxy.Read(buf[0:])
		if err != nil {
			log.Println("UNREG SERVICE")
			break
		}
	}
}

func initUserSession(conn net.Conn, chSession chan net.Conn) {
	chSession <- conn
}

//处理最终用户的连接
func userConnect(conn net.Conn, bufReaded []byte, chSession chan net.Conn) {
	if clientProxy == nil {
		conn.Write([]byte("NO SERVICE"))
		util.CloseConn(conn)
		return
	}
	_, err := clientProxy.Write([]byte(util.P2C_NEW_SESSION))
	if err != nil {
		conn.Write([]byte("SERVICE FAIL"))
		util.CloseConn(conn)
		return
	}
	connSession := recvSession(chSession) // := <-chSession
	if connSession == nil {
		util.CloseConn(conn)
		return
	}
	log.Println("Transfer...")
	go util.CopyFromTo(conn, connSession, bufReaded)
	go util.CopyFromTo(connSession, conn, nil)
}

//加入超时
func recvSession(ch chan net.Conn) net.Conn {
	var conn net.Conn = nil
	select {
	case conn = <-ch:
	case <-time.After(time.Second * 5):
		conn = nil
	}
	return conn
}
