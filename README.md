
# Exercise2
Mandatory Exercise 2

To run the application please do the following in the specified order, thank you

To run the service with three clients:

Open three terminals in the Excise2 folder.

1. In the first terminal run `go run main.go :7374 :7375`
2. In the second terminal run `go run main.go :7375 :7373`
3. In the third terminal run `go run main.go :7373 :7374` 

This starts three nodes each with one peer. Where the first port is the port number of the node, and the second port number is the port of the node's peer.

The node with the number ':7373' should always be run last to ensure that all peers are ready before the actual program starts.

Note: the system can be run with more peers as long as the above pattern is replicated.