package main

import (
	"1mconnections/golang/host"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	connected int64
	failed    int64
)

func main() {

	// create hosturl object
	hosturl := host.NewHostURL()

	// Listen for incoming connections.

	go func() {
		start := time.Now()
		for {
			fmt.Printf("server elapsed=%0.0fs connected=%d failed=%d\n", time.Since(start).Seconds(), atomic.LoadInt64(&connected), atomic.LoadInt64(&failed))
			time.Sleep(1 * time.Second)
		}
	}()

	// number of servers to create

	num := os.Args[1]
	num_servers, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}

	for n := 0; n < num_servers; n++ {
		if n%2000 == 0 {
			if err := hosturl.UpgradeHostURL(); err != nil {
				panic(err)
			}
		}
		addr := hosturl.String() + ":" + strconv.Itoa(2000+n%2000)
		go create_server(addr)
		fmt.Printf("Listening on host: %s, port: %s\n", hosturl.String(), strconv.Itoa(2000+n%2000))
	}

	for {
	}

}

func create_server(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()
	// host, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Listening on host: %s, port: %s\n", host, port)

	for {
		// Listen for an incoming connection
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		// Handle connections in a new goroutine
		go func(conn net.Conn) {
			atomic.AddInt64(&connected, 1)
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				atomic.AddInt64(&failed, 1)
				atomic.AddInt64(&connected, -1)
				return
			}
			fmt.Printf("Message received: %s\n", string(buf[:len]))
			for {
				len, err := conn.Read(buf)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					atomic.AddInt64(&connected, -1)
					fmt.Printf("Error reading: %#v\n", err)
					return
				}
				if string(buf[:len]) == "exit" {
					fmt.Println("Exiting...")
					conn.Write([]byte("Exiting...\n"))
					conn.Close()
					return
				}
				conn.Write([]byte("Message received.\n"))
			}

		}(conn)
	}

}
