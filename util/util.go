package util

import (
	"encoding/binary"
	"errors"
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
	SEPS            = "\n"
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

func WriteString(w io.Writer, str string) (int, error) {
	binary.Write(w, binary.LittleEndian, int32(len(str)))
	return w.Write([]byte(str))
}

const MAX_STRING = 10240

func ReadString(r io.Reader) (string, error) {
	var size int32
	err := binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return "", err
	}
	if size > MAX_STRING {
		return "", errors.New("too long string")
	}

	buff := make([]byte, size)
	n, err := r.Read(buff)
	if err != nil {
		return "", err
	}
	if int32(n) != size {
		return "", errors.New("invalid string size")
	}
	return string(buff), nil
}
