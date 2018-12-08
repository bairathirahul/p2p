package main

import (
	"sync"
)

// mutex : Mutex for the knownHosts map
var mutex = sync.Mutex{}
var knownHosts = make(map[string]string)

// getHosts : Get list of known hosts
func getHosts() map[string]string {
	return knownHosts
}

// addHosts : Add new hosts to the list of known hosts
func addHosts(newHosts map[string]string) {
	// Acuqire mutex lock
	mutex.Lock()
	// Add known hosts
	for host, name := range newHosts {
		knownHosts[host] = name
	}
	// Update
	mutex.Unlock()
}

// removeHost : Remove host from the list of known hosts
func removeHost(host string) {
	// Acquire mutex lock
	mutex.Lock()
	// Delete entry from known hosts dictionary
	delete(knownHosts, host)
	// Release mutex lock
	mutex.Unlock()
}
