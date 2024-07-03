package transport

// Transport is responsible for communication between client and server
// Main features to be covered are:
// 	- Listen for incoming connections
//  - Accept incoming connections
//  - Handle connections

type Transport interface {
	acceptConnection() error

	requestConnection() error

	sendMessage() error

}
