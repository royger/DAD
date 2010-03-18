package main

import(
	"flag";
	"net";
	"fmt";
	"os";
	"bufio"
	"log"
	)
	
const (
	BUFFER_SIZE = 512
	)

func echoClient(c net.Conn) {
	// Close the socket when the function returns
	defer c.Close()
	// Allocate buffer
	//b := make([]byte, BUFFER_SIZE)
	// Read every line untill we receive a EOF (^D)
	input := bufio.NewReader(os.Stdin)
	output := bufio.NewReader(c)
	for line, err := input.ReadString('\n'); err != os.EOF; line, err = input.ReadString('\n') {
		if err != nil {
			panic("bufio.ReadString: ", err.String())
		}
		// Write the buffer to the socket
		nw, err := c.Write([]byte(line))
		if err != nil {
			panic("Write: ", err.String())
		} else if len(line) != nw {
			panic("nr != nw")
		}
		// Read from the socket
		line, err := output.ReadString('\n')
		if err != nil {
			panic("bufio.ReadString: ", err.String())
		}
		// Print the string received from the socket
		fmt.Println(line)
	}
}

func main() {
	// Parse input flags
	var socket *string = flag.String("addr", "", "IP:PORT to connect to")
	flag.Parse()
	// Get the adress
	addr, err := net.ResolveTCPAddr(*socket)
	if err != nil {
		panic("ResolveTCPAddr: ", err.String())
	}
	// Create the connection
	c, err := net.DialTCP("tcp4", nil, addr)
	// Print info about the unix socket
	log.Stderr("Connected to", c.RemoteAddr())
	
	// Main bucle
	echoClient(c)
}
