package main

import (
	proto "Exercise2/grpc"
	"bufio"
	"context"
	"fmt"
	//"google.golang.org/grpc"
	"log"
	"os"
)

type Client struct {
	proto.UnimplementedCriticalSectionServiceServer
	id        int
	ipAndPort string
	timeStamp int
	queue     []string
}

func main() {
	// Get the client port in the form ip:port
	clientIpAndPort := os.Args[1]

	otherClients := getOtherClientRoutes("ipsAndPorts.txt", clientIpAndPort)
	for j := 0; j < len(otherClients); j++ {
		fmt.Println(otherClients[j])
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

func (c *Client) Request(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}

func (c *Client) Reply(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}

func (c *Client) Release(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}
