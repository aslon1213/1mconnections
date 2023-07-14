package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	// client_number := make(chan int)

	PORT := ":" + arguments[1]
	clients := []*net.Conn{}
	for i := 0; i < 200; i++ {
		fmt.Println(i)
		function := func() {

			client, err := net.Dial("tcp4", PORT)
			if err != nil {
				fmt.Println(err)
				return
			}
			clients = append(clients, &client)
			defer client.Close()

			// fmt.Println("connected to server:", client.RemoteAddr().String())
			// fmt.Println("sending request from:", client.LocalAddr().String())

			// fmt.Println("Sleeping for 25 seconds...")
			time.Sleep(25 * time.Second)

		}
		go function()
	}

	fmt.Println("clients has connected:", len(clients))

	time.Sleep(25 * time.Second)
}
