package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	TYPE = "tcp"
	HOST = "localhost"
	PORT = "8888"
)

func main() {

	//get the number of clients to connect to the server
	num := os.Args[1]
	num_2 := os.Args[2]
	fmt.Println("Number of clients:", num)
	number_of_clients, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}
	number_of_clients_2, err := strconv.Atoi(num_2)

	if err != nil {
		panic(err)
	}

	for j := 0; j < number_of_clients_2; j++ {
		for i := 0; i < number_of_clients; i++ {
			go create_client("localhost", strconv.Itoa(2000+j))
			fmt.Println("Client", i+1, "has connected to the server")
		}
	}

	for {

	}

}

func create_client(host string, port string) {

	// listen to incoming messages
	conn, err := net.Dial(TYPE, host+":"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// send to conn a message "Hello"
	_, err = conn.Write([]byte("Hello"))
	if err != nil {
		panic(err)
	}

	for {
	}

}
