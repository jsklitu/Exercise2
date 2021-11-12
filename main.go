package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Client struct {
	id        int
	ipAndPort string
	timeStamp int
	queue     []string
}

func main() {
	// Get the client port in the form ip:port
	clientIpAndPort := os.Args[1]

	clients := getOtherClientRoutes("ipsAndPorts.txt", clientIpAndPort)
	for j := 0; j < len(clients); j++ {
		fmt.Println(clients[j])
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

/*func Request(ctx context.Context, in *proto.Message) (*proto.Close, error) {

}

func Reply(ctx context.Context, in *proto.Message) (*proto.Close, error) {

}

func Release(ctx context.Context, in *proto.Message) (*proto.Close, error) {

}



func sendRequestMessage(peer *proto.Peer) {
	lamportTime++
	newMsg := &proto.Message{
		Id:        peer.id,
		Message:   "I wanna get in!!",
		Timestamp: lamportTime,
	}
	_, requestErr := peer.BroadcastMessage(context.Background(), newMsg)

	if requestErr != nil {
		log.Println("error sending join message: ", requestErr)
	}
}*/
