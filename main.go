package main

import (
	proto "Exercise2/grpc"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sync"

	"context"
	"log"
	"os"
)

type Client struct {
	proto.UnimplementedCriticalSectionServiceServer
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
		id:        clientPort,
		timeStamp: 1,
		queue:     []string{},
		peers:     []proto.CriticalSectionServiceClient{},
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

	waitForReady()

	sendARequest(client)

	go func() { // Wait for our wait group decrementing
		waiter.Wait()
		close(done)
	}()

	fmt.Println("66")

	<-done // Wait until done sends back some data

	fmt.Println("69")
}

func sendARequest(client *Client) {
	for _, peer := range client.peers {
		newMessage := &proto.Message{
			Request: "send a request",
			Id:      client.id,
		}
		returnRequest, err := peer.Request(context.Background(), newMessage)
		if err != nil {
			return
		}
		fmt.Println(returnRequest.Request, " from ", returnRequest.Id)
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

func (c *Client) Receive(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer log.Println("im done with recieve method")
	log.Println("Received a message from ", in.Id, ". The critical section key is: ", in.CriticalSection)
	if c.wantAccess {
		log.Println("Im in the critical section")
		in.CriticalSection++
		time.Sleep(3 * time.Second)
		c.wantAccess = false
		log.Println("Increased the critical section key to: ", in.CriticalSection)
		log.Println("I'm leaving the critical section")
	} else {
		fmt.Println("I don't want access right now – passing the key on")
	}

	fmt.Println("Below critical section section")

	in.Id = c.id
	_, err := c.peer.Receive(ctx, in)
	if err != nil {
		return nil, err
	}
	returnMsg := &proto.Close{}

	return returnMsg, nil
}
