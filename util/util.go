package util

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

//TOKEN
const (
	TOKEN_LEN       = 4
	C2P_CONNECT     = "C2P0"
	C2P_SESSION     = "C2P1"
	C2P_KEEP_ALIVE  = "C2P2"
	P2C_NEW_SESSION = "P2C1"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func Conn2Str(conn net.Conn) string {
	return conn.LocalAddr().String() + " <-> " + conn.RemoteAddr().String()
}

func CopyFromTo(r, w io.ReadWriteCloser, buf []byte) {
	defer CloseConn(r)
	if buf != nil && len(buf) > 0 {
		_, err := w.Write(buf)
		if err != nil {
			return
		}
	}
	io.Copy(r, w)
}

func CloseConn(a io.ReadWriteCloser) {
	fmt.Println("CLOSE")
	a.Close()
}
