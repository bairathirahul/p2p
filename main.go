package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Configuration : configuration
type Configuration struct {
	Name       string   `json:"name"`
	Port       string   `json:"port"`
	KnownHosts []string `json:"knownHosts"`
}

// Message : Incoming and outgoing message structure
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var config Configuration
var channelSend chan Message
var channelReceive chan Message

func main() {
	// Read the configuration file
	log.Println("Reading configuration file")
	configJSON, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Configuration file not found")
		return
	}
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Println("Error in configuration file - " + err.Error())
		return
	}

	// Initialize channels
	channelSend = make(chan Message)
	channelReceive = make(chan Message)

	// Initialize Server
	go initServer(channelReceive)

	// Send introduction to known hosts
	go initIntroduction()

	// Initialize heartbeats
	go initHeartbeats()

	// Initialize Web UI
	initUI(channelReceive, channelSend)

	// Shutdown heartbeats
	shutdownHeartBeats()

	// Shudown server
	shutdownServer()

	// Exit ui
	exitUI()
}
