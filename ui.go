package main

import (
	"encoding/json"
	"github.com/zserge/webview"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"time"
)

// UIBinding : Binding to the web view
type UIBinding struct {
	Name string `json:"name"`
}

var view webview.WebView

// SendChat : Send message to all hosts, this method is bound in the webview
func (binding *UIBinding) SendChat(text string) {
	// Prepare chat message
	content := make(map[string]string)
	content["name"] = config.Name
	content["created"] = strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(content["created"])
	content["text"] = text
	contentJSON, _ := json.Marshal(content)
	message := Message{"CHAT", string(contentJSON)}

	// Broadcast message to all nodes
	broadcastMessage(message)
}

// initUI : Initialize web user interface
func initUI(channelReceiver chan Message, channelSend chan Message) {
	// Read index.html file
	markup, err := ioutil.ReadFile("assets/index.html")
	if err != nil {
		log.Println("Markup file not found ", err.Error())
		return
	}

	// Initialize webview
	view = webview.New(webview.Settings{
		Title: "P2P Chat",
		URL:   "data:text/html," + url.PathEscape(string(markup)),
	})

	view.Dispatch(func() {
		// Assign binding to the webview
		binding := UIBinding{}
		binding.Name = config.Name
		view.Bind("binding", &binding)
	})

	go func() {
		for {
			message := <-channelReceiver
			log.Println("Received message for UI ", message)
			switch message.Type {
			case "INTRODUCTION":
				// Send received introduction to the webview
				view.Dispatch(func() {
					view.Eval("receivedIntroduction(" + message.Content + ")")
				})

			case "CHAT":
				// Send received chat to the webview
				view.Dispatch(func() {
					view.Eval("receivedChat(" + message.Content + ")")
				})

			case "DISCONNECTION":
				// Send disconnection notification to the webview
				view.Dispatch(func() {
					view.Eval("receivedDisconnection(" + message.Content + ")")
				})
			}
		}
	}()

	// Launch webview
	view.Run()
}

// exitUI : Close the webview and exit application
func exitUI() {
	view.Exit()
}
