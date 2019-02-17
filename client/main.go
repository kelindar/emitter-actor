package main

import (
	"fmt"
	"strings"

	emitter "github.com/emitter-io/go/v2"
)

var room = "room1"

func main() {
	client, err := emitter.Connect("", onMessageReceived)
	if err != nil {
		panic(err)
	}

	// Create a private link so we can receive dedicated replies
	if _, err := client.CreatePrivateLink("4Xjj-2MKX6TxN8LqKfNJ6cBmIZniDgLO", "actor/", "s", true); err != nil {
		panic(err)
	}

	send(client, "enter", "")
	println("enter some text or 'q' to exit:")
	for {
		var text string
		if _, err := fmt.Scanln(&text); err != nil {
			println(err.Error())
			return
		}

		if strings.ToLower(text) == "q" {
			return
		}

		send(client, text, "")
	}
}

func send(c *emitter.Client, command, message string) {
	payload := fmt.Sprintf("%s %s %s", c.ID(), command, message)

	// The key is for "actor/#/", allowed only PUBLISH and nothing else
	err := c.Publish("CJdKnIsQoMFxvSfmBqLz3LXbkdCfHbGW", topic(room), payload)
	if err != nil {
		println(err.Error())
	}
}

func onMessageReceived(c *emitter.Client, msg emitter.Message) {
	request := strings.SplitN(string(msg.Payload()), " ", 3)
	if len(request) != 3 {
		return
	}

	command := request[1]
	message := request[2]
	switch command {
	case "move":
		println("You walk through the door.", message)
		room = message
		send(c, "enter", "")
	case "tell":
		println(message)
	default:
		println("unknown message: ", msg.Topic(), string(msg.Payload()))
	}
}

// Topic creates a topic name from an actor name
func topic(name string) string {
	return fmt.Sprintf("actor/%s/", name)
}
