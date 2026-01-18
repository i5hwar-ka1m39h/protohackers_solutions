package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/textproto"
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

type Request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func validateRequest(req Request) error {

	if req.Method == nil {
		return errors.New("method not found")
	}
	if *req.Method != "isPrime" {
		return errors.New("invalid method")
	}

	if req.Number == nil {
		return errors.New("missing number")
	}

	return nil
}

func isPrime(num float64) bool {

	if num != math.Trunc(num) {
		return false
	}

	intNum := int64(num)

	if intNum < 2 {
		return false
	}

	for i := int64(2); i*i <= intNum; i++ {
		if intNum%i == 0 {
			return false
		}
	}

	return true

}

func doSomeStuff(c net.Conn, id int32) {
	fmt.Println("working for id:", id, "from endpoitn", c.RemoteAddr().String())

	defer c.Close()
	conn := textproto.NewReader(bufio.NewReader(c))

	malformed := []byte("{}\n")

	for {
		line, err := conn.ReadLineBytes()

		if err != nil {
			if err == io.EOF {
				log.Println("end of line error", err)
			}
			log.Println("read line bytes error ", err)
		}

		var req Request
		err = json.Unmarshal(line, &req)
		if err != nil {

			log.Println("json unmarshal error", err)

			c.Write(malformed)
			return
		}

		err = validateRequest(req)
		if err != nil {
			log.Println("validator error:", err)
			c.Write(malformed)
			return
		}

		resp := Response{
			Method: "isPrime",
			Prime:  isPrime(float64(*req.Number)),
		}

		respData, err := json.Marshal(&resp)
		if err != nil {
			log.Fatalln("error in marshal", err)
		}

		respData = append(respData, []byte("\n")...)
		c.Write(respData)

	}

}
