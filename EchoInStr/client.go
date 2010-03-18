package main

import(
	"flag";
	"net";
	"log";
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
