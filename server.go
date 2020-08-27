package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// Messages Delimeter
const (
	QuitSign      = "/quit!"
)

// ScanCRLF is a
func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte(QuitSign)); i >= 0 {
		return i + 7, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func main() {
	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)
	for {
	reader := bufio.NewReader(conn)

	ba, _, err := reader.ReadLine()
	log.Println(string(ba))
	if err != nil {
		if err == io.EOF {
			fmt.Println(err)
		}
		break
	}

	bf, er := strconv.Atoi(string(ba))
	if er != nil {
		log.Println(er)
	}
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, bf+len(QuitSign))
	scanner.Buffer(buf, bf+len(QuitSign))
	scanner.Split(ScanCRLF)
	
		ok := scanner.Scan()

		if !ok {
			break
		}
		log.Println(len(scanner.Text()))
		handleMessage(scanner.Text(), conn)
	}
	// log.Println(bf)
	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println(">" + message)

}
