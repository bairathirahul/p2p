package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)

// stopHearbeat : Flag to stop the heartbeat thread
var stopHeartbeat = false

// initIntroduction : Send introduction message to all known hosts
func initIntroduction() {
	log.Println("Sending introductions to ", config.KnownHosts)
	if len(config.KnownHosts) > 0 {
		contentMap := map[string]string{
			"name": config.Name,
			"port": config.Port,
		}
		contentJSON, _ := json.Marshal(contentMap)
		message := Message{"INTRODUCTION", string(contentJSON)}
		messageJSON, _ := json.Marshal(message)
		messageJSON = append(messageJSON, 10)
		for _, host := range config.KnownHosts {
			go sendIntroduction(messageJSON, host)
		}
	}
}

// sendIntroduction : Send the introduction message to specified host
func sendIntroduction(messageJSON []byte, host string) {
	// Connect to the host server
	client, err := net.Dial("tcp", host)
	if err != nil {
		log.Println("Cannot connect to the client " + host)
		log.Println(err.Error())
		return
	}

	// Send message to the host
	n, err := client.Write(messageJSON)
	if err != nil {
		log.Println("Connection write error " + err.Error())
		return
	}
	log.Println(strconv.Itoa(n) + " bytes sent to " + host)

	// Read connection data
	data, err := bufio.NewReader(client).ReadBytes(10)
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

	log.Println("Received introduction response", message)
	if message.Type == "HOSTS" {
		// Unmarshal the response
		var contentMap map[string]string
		var newHosts map[string]string
		err = json.Unmarshal([]byte(message.Content), &contentMap)
		if err != nil {
			log.Println("Error unmarshalling response", err.Error())
			return
		}
		err = json.Unmarshal([]byte(contentMap["hosts"]), &newHosts)
		if err != nil {
			log.Println("Error unmarshalling response", err.Error())
		}

		// Send introduction to any new host
		for newHost := range newHosts {
			if _, ok := knownHosts[newHost]; !ok {
				go sendIntroduction(messageJSON, newHost)
			}
		}

		// Add host to the known hosts list
		newHosts[host] = contentMap["name"]
		addHosts(newHosts)
	}
}

// broadcastMessage : Broadcast chat message to all known hosts
func broadcastMessage(message Message) {
	messageJSON, _ := json.Marshal(message)
	messageJSON = append(messageJSON, 10)
	log.Println("Received message for broadcasting", message)
	log.Println("Sending message to", knownHosts)
	for host := range getHosts() {
		go func(messageJSON []byte, host string) {
			// Connect to the host server
			client, err := net.Dial("tcp", host)
			if err != nil {
				log.Println("Cannot connect to the client " + host)
				log.Println(err.Error())
				return
			}

			// Send message to the host
			n, err := client.Write(messageJSON)
			if err != nil {
				log.Println("Connection write error " + err.Error())
				return
			}
			log.Println(strconv.Itoa(n) + " bytes sent to " + host)
		}(messageJSON, host)
	}
}

// initHeartbeats : Initialize the thread to send heartbeat to all hosts every 5 seconds
func initHeartbeats() {
	log.Println("Initializing heartbeats")
	contentMap := map[string]string{
		"name": config.Name,
		"port": config.Port,
	}
	contentJSON, _ := json.Marshal(contentMap)
	message := Message{"HEARTBEAT", string(contentJSON)}
	messageJSON, _ := json.Marshal(message)
	for !stopHeartbeat {
		time.Sleep(5 * time.Second)
		log.Println("Sending heartbeats to ", getHosts())
		for host := range getHosts() {
			if host == "self" {
				continue
			}
			go func(messageJSON []byte, host string) {
				// Connect to the host server
				client, err := net.Dial("tcp", host)
				if err != nil {
					// Remove host on error
					log.Println("Cannot connect to the client " + host)
					log.Println(err.Error())

					// Send message of disconnection
					content := map[string]string{"name": knownHosts[host]}
					contentJSON, _ := json.Marshal(content)
					message := Message{"DISCONNECTION", string(contentJSON)}
					channelReceive <- message

					// Remove host
					removeHost(host)
					return
				}

				// Send message to the host
				n, err := client.Write(messageJSON)
				if err != nil {
					// Remove host on error
					log.Println("Connection write error " + err.Error())

					// Send message of disconnection
					content := map[string]string{"name": knownHosts[host]}
					contentJSON, _ := json.Marshal(content)
					message := Message{"DISCONNECTION", string(contentJSON)}
					channelReceive <- message

					// Remove host
					removeHost(host)
					return
				}
				log.Println(strconv.Itoa(n) + " bytes sent to " + host)
			}(messageJSON, host)
		}
	}
	log.Println("Heartbeats stopped")
}

// shutdownHeartBeats : Set the flag to stop heart beats
func shutdownHeartBeats() {
	log.Println("Stopping heartbeats")
	stopHeartbeat = true
}
