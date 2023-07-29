package main

import (
	"1mconnections/golang/host"
	"fmt"
	"log"
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

	host := host.NewHostURL()

	//get the number of clients to connect to the server
	num := os.Args[1]
	fmt.Println("Number of clients:", num)
	number_of_clients, err := strconv.Atoi(num)

	if err != nil {
		panic(err)
	}

	// choose the correct

	for i := 0; i < number_of_clients; i++ {

		if i%2000 == 0 {
			if err := host.UpgradeHostURL(); err != nil {
				panic(err)
			}
		}

		go create_client(host.String(), strconv.Itoa(2000+i%2000), 49152+i%2000)
		fmt.Println("Connected to " + host.String() + ":" + strconv.Itoa(2000+i%8000) + " from " + host.String() + ":" + strconv.Itoa(49152+i%16000))

	}
	for {

	}

}

func create_client(host string, remoteport string, localport int) {

	// listen to incoming messages with
	dialer := net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(HOST),
			Port: localport,
		},
	}
	conn, err := dialer.Dial(TYPE, host+":"+remoteport)

	if err != nil {
		log.Println(err)

		return
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
