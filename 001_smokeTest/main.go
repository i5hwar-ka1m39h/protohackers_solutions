package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listner, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalln("error occured while creating the listener", err)
	}

	fmt.Println("server is listening on port 3001")
	defer listner.Close()

	var id int32 = 0
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatalln("error accepting the connection", conn)
		}

		id++
		go doSomeStuff(conn, id)

	}
}

func doSomeStuff(c net.Conn, id int32) {
	var buff bytes.Buffer
	fmt.Println("process running for id :", id)
	fmt.Println("remote address :", c.RemoteAddr().String())

	readBytes, err := io.Copy(&buff, c)
	if err != nil {
		log.Fatalln("error occured while reading the bytes", err)
	}

	fmt.Println("read bytes are :", readBytes)

	bytesWritten, err := c.Write(buff.Bytes())
	if err != nil {
		log.Fatalln("error occured while wrinting buffere:", err)
	}
	fmt.Println("written bytes are :", bytesWritten)

	defer c.Close()
}
