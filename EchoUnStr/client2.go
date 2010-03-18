package main

import(
	"flag";
	"net";
	"fmt";
	"os";
	)
	
const (
	BUFFER_SIZE = 512
	)

func echoClient(c net.Conn) {
	// Close the socket when the function returns
	defer c.Close()
	// Allocate buffer
	b := make([]byte, BUFFER_SIZE)
	// Read every line untill we receive a EOF (^D)
	for nr, err := os.Stdin.Read(b); err != os.EOF; nr, err = os.Stdin.Read(b) {
		if err != nil {
			panic("Read: ", err.String())
		}
		// Write the buffer to the socket
		nw, err := c.Write(b[0:nr])
		if err != nil {
			panic("Write: ", err.String())
		} else if nr != nw {
			panic("nr != nw")
		}
		// Read from the socket
		nr, err := c.Read(b)
		if err != nil {
			panic("Read: ", err.String())
		}
		// Print the string received from the socket
		fmt.Println(string(b[0:nr]))
	}
}

func main() {
	// Parse input flags
	var socket *string = flag.String("socket", "", "Unix domain socket to connect to")
	flag.Parse()
	// Get the adress
	addr, err := net.ResolveUnixAddr("unix", *socket)
	if err != nil {
		panic("ResolveUnixAddr: ", err.String())
	}
	// Create the connection
	c, err := net.DialUnix("unix", nil, addr)
	// Print info about the unix socket
	fmt.Println("Connected to", c.RemoteAddr())
	
	// Main bucle
	echoClient(c)
}
