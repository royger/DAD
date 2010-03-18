package main

import(
	"flag";
	"net";
	"fmt";
	"os";
	"runtime";
	)

func createConn(socket, domain string) (c *net.UnixConn){
	addr, err := net.ResolveUnixAddr(domain, socket)
	if err != nil {
		panic("ResolveUnixAddr: ", err.String())
	}
	// Create the connection
	Connection: c, err = net.DialUnix(domain, nil, addr)
	if err != nil {
		if err.(*net.OpError).Error == os.ECONNREFUSED {
			goto Connection
		}
		panic("DialUnix: ", err.String())
	}
	return
}

func createTestClient(socket, domain string, sem chan bool) {
	c := createConn(socket, domain)
	defer c.Close()
	// Message to transmit
	b := "Això és una prova"
	// Write message to socket
	nw, err := c.Write([]byte(b))
	if err != nil {
		panic("Write: ", err.String())
	} else if len(b) != nw {
		panic("nr != nw")
	}
	// Read from the socket
	nr, err := c.Read([]byte(b))
	if err != nil {
		panic("Read: ", err.String())
	} else if nr != nw {
		panic("nr != nw")
	}
	// Increase the semaphore
	sem <- true
}

func main() {
	// Parse input flags
	var socket *string = flag.String("socket", "", "Unix domain socket to connect to")
	var num_clients *int = flag.Int("num_clients", 100, "Number of clients to launch")
	var concurrency *int = flag.Int("concurrency", 10, "Number of clients running concurrently")
	var num_cpu *int = flag.Int("cpu_use", 2, "Number of CPUs to use")
	flag.Parse()
	// Set thew number of CPUs to use
	runtime.GOMAXPROCS(*num_cpu)
	// Create a chan to comunicate with the processes
	// This channel will imitate a semaphore
	p := 0
	sem := make(chan bool, *concurrency)
	// Place the initial values in the semaphore
	for i := 0; i < *concurrency; i++ {
		sem <- true
	}
	// Run the test clients
	for <- sem ;p < *num_clients; _, p = <- sem, p+1 {
		go createTestClient(*socket, "unix", sem)
	}
	// Wait for the remaining processes to finish
	for j :=  0; j < *concurrency-1; _, j = <- sem, j+1 { }
	fmt.Println(p, "processes finished")
}
