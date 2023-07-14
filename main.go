package main

import (
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

type HostURL struct {
	First_Chunk  int // 127
	Second_Chunk int // 0-255
	Third_Chunk  int // 0-255
	Fourth_Chunk int // 0-255
}

func NewHostURL() *HostURL {
	return &HostURL{
		First_Chunk:  127,
		Second_Chunk: 0,
		Third_Chunk:  0,
		Fourth_Chunk: 0,
	}
}

func (h *HostURL) UpgradeHostURL() error {

	if h.Fourth_Chunk < 255 {
		h.Fourth_Chunk++
		return nil
	} else if h.Third_Chunk < 255 {
		h.Third_Chunk++
		h.Fourth_Chunk = 1
		return nil
	} else if h.Second_Chunk < 255 {
		h.Second_Chunk++
		h.Third_Chunk = 0
		h.Fourth_Chunk = 0
		return nil
	} else {
		return fmt.Errorf("HostURL is out of range")
	}

}

func (h *HostURL) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", h.First_Chunk, h.Second_Chunk, h.Third_Chunk, h.Fourth_Chunk)
}

func main() {

	// create hosturl object
	hosturl := NewHostURL()

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
	for n := 0; n < num_servers/8000; n++ {
		if err := hosturl.UpgradeHostURL(); err != nil {
			panic(err)
		}
		j := 0
		for j < 8000 {
			addr := hosturl.String() + ":" + strconv.Itoa(2000+j)
			go create_server(addr)
			j++
		}
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
	host, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on host: %s, port: %s\n", host, port)

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
