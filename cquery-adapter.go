package main

import (
	"flag"
	"net"
	"os"
	"runtime"
)

func readStdinSendConn(tcpConn net.Conn, logPath *string) {
	var buffer [6 * 1024 * 1024]byte
	var flog *os.File
	if logPath != nil || len(*logPath) == 0 {
		flog, _ := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer flog.Close()
	}

	for {
		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			panic(err)
		}

		n, err = tcpConn.Write(buffer[:n])
		if err != nil {
			panic(err)
		}

		if flog != nil {
			flog.Write(buffer[:n])
		}
	}
}

func recvConnWriteStdout(tcpConn net.Conn, logPath *string) {
	var buffer [6 * 1024 * 1024]byte
	var flog *os.File
	if logPath != nil || len(*logPath) == 0 {
		flog, _ := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		defer flog.Close()
	}

	for {
		n, err := tcpConn.Read(buffer[:len(buffer)-1])
		if err != nil {
			panic(err)
		}

		n, err = os.Stdout.Write(buffer[:n])
		if err != nil {
			panic(err)
		}

		if flog != nil {
			flog.Write(buffer[:n])
		}
	}
}

func main() {
	inLogPath := flag.String("inlog", "", "input message log file path")
	outLogPath := flag.String("outlog", "", "output message log file path")
	address := flag.String("h", "", "ip address for cquery server [ip:hsot]")
	flag.Bool("language-server", true, "cquery default arg")
	flag.Parse()

	tcpConn, err := net.Dial("tcp", *address)
	if err != nil {
		panic(err)
	}
	defer tcpConn.Close()

	go readStdinSendConn(tcpConn, inLogPath)
	go recvConnWriteStdout(tcpConn, outLogPath)

	for {
		runtime.Gosched()
	}
}
