package main

import (
	proto "Exercise2/grpc"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"

	"context"
	"log"
	"os"
)

type Connection struct {
	stream proto.CriticalSectionService_ReceiveClient
	error  chan error
}

type Client struct {
	proto.UnimplementedCriticalSectionServiceServer
	stream     proto.CriticalSectionService_ReceiveServer
	mu         sync.Mutex
	id         string
	peer       proto.CriticalSectionServiceClient
	peerId     string
	wantAccess bool
}

func main() {
	done := make(chan int)
	waiter := &sync.WaitGroup{}

	clientPort := os.Args[1]
	peerPort := os.Args[2]
	if len(clientPort) > 5 || len(peerPort) > 5 {
		log.Fatalf("Invalid input")
	}

	client := &Client{
		id:         clientPort,
		wantAccess: false,
		peerId:     peerPort,
	}

	clientServerGrpc := grpc.NewServer()           // Start server
	listener, err := net.Listen("tcp", clientPort) // Listen at the client's port

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Println("Started server at port: ", clientPort)

	proto.RegisterCriticalSectionServiceServer(clientServerGrpc, client)

	go serverRunning(clientServerGrpc, listener, waiter)
	connectPeer(client, peerPort)

	if client.id == ":7373" {
		go startCircle(client)
	}

	go requestAccess(client)

	go func() { // Wait for our wait group decrementing
		waiter.Wait()
		close(done)
	}()

	<-done // Wait until done sends back some data

	fmt.Println("69")
}

func startCircle(client *Client) {
	log.Println("starting circle")
	stream, err := client.peer.Receive(context.Background())
	if err != nil {
		log.Println("Could not start circle")
	}
}

func requestAccess(client *Client) {
	for {
		if client.wantAccess != true {
			time.Sleep(10 * time.Second)
			client.mu.Lock()
			client.wantAccess = true
			log.Println("i have changed to true")
			client.mu.Unlock()
		}
	}
}

func serverRunning(clientServerGrpc *grpc.Server, listener net.Listener, wait *sync.WaitGroup) {
	wait.Add(1)
	serveError := clientServerGrpc.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

// Connect with the neighboring peer
func connectPeer(client *Client, peerPort string) {
	conn, err := grpc.Dial("localhost"+peerPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Client could not connect to peer " + peerPort)
	}
	newPeer := proto.NewCriticalSectionServiceClient(conn)
	client.peer = newPeer
}

func (c *Client) Receive(str proto.CriticalSectionService_ReceiveServer) error {
	if c.peerStream.stream == nil {
		stream, err := c.peer.Receive(context.Background())
		if err != nil {
			return err
		}
		c.peerStream.stream = stream
		c.peerStream.error = make(chan error)
	}
	c.stream = str
	go waitForMessage(c)

	return <-c.peerStream.error
}
