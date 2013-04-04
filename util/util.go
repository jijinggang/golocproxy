package util

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func Conn2Str(conn net.Conn) string {
	return conn.LocalAddr().String() + " <-> " + conn.RemoteAddr().String()
}

const (
	TOKEN_LEN       = 4
	C2P_CONNECT     = "C2P0"
	C2P_SESSION     = "C2P1"
	P2C_NEW_SESSION = "P2C1"
)
