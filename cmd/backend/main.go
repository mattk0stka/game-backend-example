package main

import (
	"bufio"
	"io"
	"log"
	"net"
)

/*
there are multiple approaches to distinguish different FlatBuffer tables, send
over the network, but the most recommended and idiomatic way (endorsed by the
FlatBuffer specification itself) is: using a `union` inside a wrapper/root table
*/
func main() {
	/*
		- if the function returns successfully, the listener is bound to the specified IP address and port.
		Binding: the operating system has exclusively assigned the port on the given IP address to the listener.
		OS allows not other processes to listen for incoming traffic on bound ports.

		- if port is zero or empty, a randomly port number will be assigned to the listener.
		- retrieve listener's address by calling its `Addr` method
	*/
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	log.Println("Listening on :9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	packet := make([]byte, 1024)
	tmp := make([]byte, 1024)

	for {
		_, err := conn.Read(packet)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		packet = append(packet, tmp...)
	}

	// pkt received
	bufio.NewReader(conn)
}
