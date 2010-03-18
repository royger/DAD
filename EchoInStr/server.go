package main

import(
	"net";
	"strconv";
	"io";
	"runtime";
	"flag";
	"os/signal"
	"strings"
	"log"
	)
	
func signalListener(end chan bool) {
	for {
		if n := <- signal.Incoming; strings.HasPrefix(n.String(), "SIGINT") {
			log.Stderr("Received SIGINT, ending execution")
			end <- true
			return
		} else {
			log.Stderr("Received singal", n, "ignoring")
		}
	}
}

func connListener(conn chan net.Conn, listener *net.TCPListener) {
	for {
		c, err := listener.Accept()
		if err != nil {
			panic("Accept: ", err.String())
		}
		conn <- c
	}
}
func echoServer(c net.Conn) {
	// Close the socket at the end of the comunication
	defer c.Close()
	// Copy from socket to socket
	_, err := io.Copy(c, c)
	if err != nil {
		panic("Copy: ", err.String())
	}
	log.Stderr("Finishing client")
}

func main() {
	var num_cpu *int = flag.Int("cpu_use", 1, "Number of CPUs to use")
	var port *int = flag.Int("port", 5656, "Port to listen for connections")
	flag.Parse()
	// Sets max number of CPUs to use
	runtime.GOMAXPROCS(*num_cpu)
	// Get the adress
	addr, err := net.ResolveTCPAddr(":" + strconv.Itoa(*port))
	if err != nil {
		panic("ResolveTCPAddr: ", err.String())
	}
	// Create the listener
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic("ListenTCP: ", err.String())
	}
	// Print info about the unix socket
	log.Stderr("Listening on", listener.Addr(), "using", *num_cpu, "CPUs")
	
	// Chan that signal the end of execution
	end := make(chan bool)
	// Chan to transmit new connections
	conn := make(chan net.Conn)
	// Handlers for connections and signals
	go signalListener(end)
	go connListener(conn, listener)
	
	for {
		// Wait for new connections or end of program
		select {
			case <- end:
				// End of program
				return
			case c := <- conn:
				// Launch thread
				go echoServer(c)
		}
	}
}
