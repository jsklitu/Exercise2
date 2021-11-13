package main

import (
	proto "Exercise2/grpc"
	"bufio"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sync"

	"context"
	"log"
	"os"
)

var wait *sync.WaitGroup

// The wait group will be synced as the first step
func init() {
	wait = &sync.WaitGroup{}
}

type Client struct {
	proto.UnimplementedCriticalSectionServiceServer
	id        string
	timeStamp int
	queue     []string
	peers     []proto.CriticalSectionServiceClient
}

func main() {
	done := make(chan int)
	waiter := &sync.WaitGroup{}

	// Get the client port in the form ip:port
	clientPort := os.Args[1]
	client := &Client{
		id:        clientPort,
		timeStamp: 1,
		queue:     []string{},
		peers:     []proto.CriticalSectionServiceClient{},
	}

	otherClients := getOtherClientRoutes("ports.txt", clientPort)
	fmt.Println(otherClients) // remove

	clientServerGrpc := grpc.NewServer()           // Start server
	listener, err := net.Listen("tcp", clientPort) // Listen at the client's port

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Println("Started server at port: ", clientPort)

	proto.RegisterCriticalSectionServiceServer(clientServerGrpc, client)

	go serverRunning(clientServerGrpc, listener, waiter)
	connectWithPeers(client, otherClients)

	waitForReady()

	sendARequest(client)

	go func() { // Wait for our wait group decrementing
		waiter.Wait()
		close(done)
	}()

	<-done // Wait until done sends back some data
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

func connectWithPeers(client *Client, otherClients []string) {
	for _, otherClient := range otherClients {
		conn, err := grpc.Dial("localhost"+otherClient, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Client could not connect to peer " + otherClient)
		}
		peer := proto.NewCriticalSectionServiceClient(conn)
		client.peers = append(client.peers, peer)
	}
}

func waitForReady() {
	scanner := bufio.NewScanner(os.Stdin) // Scan the input from the user through the command line
	for scanner.Scan() {
		if scanner.Text() == "ready" {
			return
		}
	}
}

func getOtherClientRoutes(fileName string, clientIpAndPort string) []string {
	// Open file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Could not read file ", err)
	}
	// Close the file at the end of the program
	defer f.Close()

	// Read file line by line
	scanner := bufio.NewScanner(f)

	var routes []string
	for scanner.Scan() {
		ipAndPort := scanner.Text()
		// Do not include the client's own port
		if ipAndPort != clientIpAndPort {
			routes = append(routes, ipAndPort)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return routes
}

func (c *Client) Request(ctx context.Context, in *proto.Message) (*proto.Message, error) {
	fmt.Println("I have received a request")
	return &proto.Message{
		Request: "Hej",
		Id:      c.id,
	}, nil
}

func (c *Client) Release(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}
