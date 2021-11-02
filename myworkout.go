package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:31564")
	checkError(err)
	req, err := bufio.NewReader(os.Stdin).ReadString('\n')
	checkError(err)
	io.WriteString(conn, req)

	msg, err := bufio.NewReader(conn).ReadString('\n')
	checkError(err)
	fmt.Println(msg)

}
