package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
		Fourth_Chunk: 1,
	}
}

func (h *HostURL) UpgradeHostURL() error {
	if h.Fourth_Chunk < 255 {
		h.Fourth_Chunk++
		return nil
	} else if h.Third_Chunk < 255 {
		h.Third_Chunk++
		h.Fourth_Chunk = 0
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
	go func() {
		start := time.Now()
		for {
			fmt.Printf("server elapsed=%0.0fs connected=%d failed=%d\n", time.Since(start).Seconds(), atomic.LoadInt64(&connected), atomic.LoadInt64(&failed))
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&connected, 1)

		// launch a new goroutine so that this function can return and the http server can free up
		// buffers associated with this connection
		go handleConnection(r, w)
	})
	hostURL := NewHostURL()
	if atomic.LoadInt64(&connected)%1000 == 0 {
		for i := 0; i < 1000; i++ {
			i := i
			fmt.Printf("Listening on %s:%d\n", hostURL.String(), 10000+i)
			go func() {
				if err := http.ListenAndServe(":"+strconv.Itoa(10000+i), nil); err != nil {
					log.Fatal(err)
				}
			}()
		}
	}

	for {

	}

}

func handleConnection(r *http.Request, w http.ResponseWriter) {

	// defer w.Close()
	for {
		fmt.Fprintf(w, time.Now().Format("15:04:05\n"))
		time.Sleep(1 * time.Second)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
