package main

import (
	"encoding/binary"
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

type Entry struct {
	timestamp int32
	val       int32
}

func doSomeStuff(c net.Conn, id int32) {
    defer c.Close()

	fmt.Println("connection for id:", id, " remote address with:", c.RemoteAddr().String())

    var store []Entry
    buf := make([]byte, 9)

    for {
        _, err := io.ReadFull(c, buf)
        if err != nil {
            return
        }

        typ := buf[0]
        a := binary.BigEndian.Uint32(buf[1:5])
        b := binary.BigEndian.Uint32(buf[5:9])

        if typ == 'I' {
            store = insertInStore(store, int32(a), int32(b))
			fmt.Println("store after ther insert", store)
            continue
        }

        if typ == 'Q' {
            res := getTheValue(store,int32(a), int32(b))
            c.Write(binary.BigEndian.AppendUint32(nil, uint32(res)))
        }
    }
}


func insertInStore(store []Entry, timeStamp int32, val int32) []Entry {
	for i, entry := range store {
		if entry.timestamp == timeStamp {
			return store
		}
		if timeStamp < entry.timestamp {
			store = append(store[:i+1], store[i:]...)
			store[i] = Entry{timestamp: timeStamp, val: val}
			return store
		}
	}
	return append(store, Entry{timestamp: timeStamp, val: val})
}

func getTheValue(store []Entry, min int32, max int32) int32 {
	if min > max {
		return 0
	}

	mean := float64(0)
	count := float64(0)

	for _, v := range store {
		if v.timestamp > max {
			break
		}
		if v.timestamp >= min {
			count++
			mean += (float64(v.val) - mean) / count
		}
	}

	if count == 0 {
		return 0
	}

	return int32(mean)
}

