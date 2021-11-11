package main

import (
	//"context"
	"fmt" //så er vi klar til debug lol
	"log"
	"os"
	"strconv"
	"time"

	//proto "Exercise2/proto"

	serf "github.com/hashicorp/serf/serf"
	//"google.golang.org/grpc/profiling/proto"
)

func main() {
	advertisePort := "5001"
	bindPort := os.Args[1]
	fmt.Println(advertisePort)
	fmt.Println(bindPort)
	cluster, err := setupCluster(
		os.Getenv("ADVERTISE_ADDR"),
		os.Getenv("CLUSTER_ADDR"),
		advertisePort,
		bindPort)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(3200 * time.Second)
	defer cluster.Leave()
}

func setupCluster(advertiseAddr string, bindAdd string, advertisePort string, bindPort string) (*serf.Serf, error) {
	conf := serf.DefaultConfig()
	conf.Init()
	conf.MemberlistConfig.AdvertiseAddr = "127.0.0.1"
	conf.MemberlistConfig.Name = bindPort

	bindPortInt, err := strconv.Atoi(bindPort)
	if err != nil {
		log.Fatal(err)
	}

	advertisePortInt, err := strconv.Atoi(advertisePort)
	if err != nil {
		log.Fatal(err)
	}

	conf.MemberlistConfig.AdvertisePort = advertisePortInt
	conf.MemberlistConfig.BindPort = bindPortInt

	// Try to create a cluster
	cluster, error := serf.Create(conf)
	if error != nil {
		log.Printf("Could not create a cluster :(")
	}

	newBind := bindPortInt - 1
	fmt.Println(bindPortInt, newBind)

	_, error = cluster.Join([]string{"127.0.0.1:" + strconv.Itoa(newBind)}, false)

	if error != nil {
		log.Printf("Could not join existing cluster - so creating own :)")
	} else {
		fmt.Println("I have joined")
	}

	return cluster, nil
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
