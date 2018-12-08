package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"strings"
)

var listener net.Listener

// initServer : Initialize and launch the server
func initServer(channelReceive chan Message) {
	// Initialize the server listener
	log.Println("Initializing server at port " + config.Port)
	listener, err := net.Listen("tcp", "127.0.01:"+config.Port)
	if err != nil {
		log.Println("Cannot initialize the server: ", err.Error())
		return
	}

	// Loop and wait for connections
	for {
		conn, err := listener.Accept()
		if err == nil {
			go handleConnection(conn, channelReceive)
		}
	}
}

// handleConnection : Handle incoming connection
func handleConnection(conn net.Conn, channelReceive chan Message) {
	remoteAddr := conn.RemoteAddr().String()
	log.Println("Received incoming message from " + remoteAddr)

	// Read connection data
	data, err := bufio.NewReader(conn).ReadBytes(10)
	if err != nil {
		log.Println("Unable to read data: " + err.Error())
		return
	}

	// Unmarshal message
	message := Message{}
	err = json.Unmarshal(data, &message)
	if err != nil {
		log.Println("Invalid data received " + err.Error())
		return
	}
	log.Println("Received message ", message)

	switch message.Type {
	case "INTRODUCTION":
		// Extract remote host information
		var contentMap map[string]string
		_ = json.Unmarshal([]byte(message.Content), &contentMap)
		remoteHost := strings.Split(remoteAddr, ":")[0] + ":" + contentMap["port"]
		userName := contentMap["name"]

		// Write message to the receiving channel
		channelReceive <- message

		// Respond with a list of known hosts
		knownHostsJSON, _ := json.Marshal(getHosts())
		contentMap = map[string]string{
			"name":  config.Name,
			"hosts": string(knownHostsJSON),
		}
		contentJSON, _ := json.Marshal(contentMap)
		message := Message{"HOSTS", string(contentJSON)}
		messageJSON, _ := json.Marshal(message)
		messageJSON = append(messageJSON, 10)
		n, err := conn.Write(messageJSON)
		if err != nil {
			log.Println("Connection write error " + err.Error())
			return
		}
		log.Println(strconv.Itoa(n) + " bytes sent")

		// Add host to the list of known hosts
		newHost := make(map[string]string)
		newHost[remoteHost] = userName
		log.Println("Adding hosts ", newHost)
		addHosts(newHost)

	case "CHAT":
		// Write message to the receiving channel
		channelReceive <- message

	case "HEARTBEAT":
		// Check for any new heartbea
	}
}

// shutdownServer : Shutdown the listener
func shutdownServer() {
	log.Println("Shutting down server")
	if listener != nil {
		listener.Close()
	}
}
