package main

import(
	"flag";
	"net";
	"fmt";
	"os";
	"io";
	)

func echoClient(c net.Conn) {
	// Close the socket when the function returns
	defer c.Close()
	// Copy from Stdin to the socket (Untill EOF is received)
	nr, err := io.Copy(c, os.Stdin)
	if err != nil {
		panic("Copy os.Stdio -> c: ", err.String())
	}
	// Copy from the socket to Stdout
	nw, err := io.Copyn(os.Stdout, c, nr)
	if err != nil {
		panic("Copy c -> os.Stdout: ", err.String())
	} else if nr != nw {
		panic("nr != nw")
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
