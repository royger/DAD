package main

import(
	"net";
	"fmt";
	"os";
	"strconv";
	"io";
	"runtime";
	"flag";
	"os/signal"
	"strings"
	)
	
func signalListener(end chan bool) {
	for {
		if n := <- signal.Incoming; strings.HasPrefix(n.String(), "SIGINT") {
			fmt.Println("Received SIGINT, ending execution")
			end <- true
			return
		} else {
			fmt.Println("Received singal", n, "ignoring")
		}
	}
}

func connListener(conn chan net.Conn, listener *net.UnixListener) {
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
	fmt.Println("Finishing client")
}

func main() {
	var num_cpu *int = flag.Int("cpu_use", 2, "Number of CPUs to use")
	flag.Parse()
	// Sets max number of CPUs to use
	runtime.GOMAXPROCS(*num_cpu)
	// Get the adress
	addr, err := net.ResolveUnixAddr("unix", "/tmp/str." + strconv.Itoa(os.Getpid()))
	if err != nil {
		panic("ResolveUnixAddr: ", err.String())
	}
	// Create the listener
	listener, err := net.ListenUnix("unix", addr)
	if err != nil {
		panic("ListenUnix: ", err.String())
	}
	// Print info about the unix socket
	fmt.Println("Listening on", listener.Addr(), "using", *num_cpu, "CPUs")
	
	// Defer deletion of UnixSocket
	defer os.Remove("/tmp/str." + strconv.Itoa(os.Getpid()))
	
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
