package main

import (
	"../util"
	"flag"
	//	"io"
	"log"
	"net"
	"strings"
	"time"
)

var (
	_port     = flag.String("p", "8010", "The Listen port of golocproxy, golocproxy client will access the port.")
	_userport = flag.String("up", "8020", "The Listen port of user connect.")
	_pwd      = flag.String("pwd", "jjg", "Password to valid Client Proxy")
)

type OnConnectFunc func(net.Conn, chan net.Conn)

func main() {
	flag.Usage = util.Usage
	flag.Parse()
	//flag.Usage()

	chSession := make(chan net.Conn, 100)
	if nil != listen(*_port, chSession, onClientConnect) {
		return
	}
	if nil != listen(*_userport, chSession, onUserConnect) {
		return
	}
	time.Sleep(999999 * time.Hour)
}

func listen(port string, chSession chan net.Conn, onConnect OnConnectFunc) error {
	server, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		log.Fatal("CAN'T LISTEN: ", err)
		return err
	}
	log.Println("listen port:", port)
	go func() {
		defer server.Close()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Println("Can't Accept: ", err)
				continue
			}
			go onConnect(conn, chSession)
		}
	}()
	return nil
}

func onClientConnect(conn net.Conn, chSession chan net.Conn) {
	strConn := util.Conn2Str(conn)
	log.Println("Proxy Client Connect:", strConn)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	msg, err := util.ReadString(conn)
	conn.SetReadDeadline(time.Time{})
	//println("Read:", string(buf[0:n]))
	if err != nil {
		log.Println("Can't Read: ", err)
		conn.Close()
		return
	}
	msgs := strings.Split(msg, "\n")
	pwd := msgs[0]
	if *_pwd != pwd {
		util.CloseConn(conn)
		return
	}
	token := msgs[1]
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

//代理客户端连接
var _clientProxy net.Conn = nil

//处理golocproxy client的连接
func clientConnect(conn net.Conn) {
	defer util.CloseConn(conn) // conn.Close()
	if _clientProxy != nil {
		util.WriteString(conn, "SERVICE EXIST")
		util.CloseConn(conn)
		return
	}
	println("REG SERVICE")
	_clientProxy = conn
	defer func() {
		_clientProxy = nil
	}()
	for {
		_, err := util.ReadString(_clientProxy)
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
func onUserConnect(conn net.Conn, chSession chan net.Conn) {
	if _clientProxy == nil {
		conn.Write([]byte("NO SERVICE"))
		util.CloseConn(conn)
		return
	}
	_, err := util.WriteString(_clientProxy, util.P2C_NEW_SESSION)
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
	go util.CopyFromTo(conn, connSession, nil)
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
