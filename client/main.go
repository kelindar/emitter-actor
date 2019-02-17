package main

import (
	"fmt"
	"strings"

	emitter "github.com/emitter-io/go/v2"
	"github.com/kelindar/emitter-actor/actor"
)

func main() {
	client, err := emitter.Connect("", func(c *emitter.Client, msg emitter.Message) {
		println("unknown message: ", msg.Topic(), string(msg.Payload()))
	})
	if err != nil {
		panic(err)
	}
	

	subKey := "4Xjj-2MKX6TxN8LqKfNJ6cBmIZniDgLO" // The key for actor/ with SUBSCRIBE and EXTEND permissions
	pubKey := "CJdKnIsQoMFxvSfmBqLz3LXbkdCfHbGW" // The key for "actor/#/" with only PUBLISH permissions

	// Create a room actor
	room := actor.Remote(pubKey, "room1", client.ID(), client)
	room.Send("enter", "")

	// Create our own actor
	self, _ := actor.New(subKey, pubKey, client.ID(), client, true)
	self.On("move", func(to, from actor.Sender, message string) {
		println("You walk through the door.", message)
		room = actor.Remote(pubKey, message, client.ID(), client)
		room.Send("enter", "")
	})

	self.On("tell", func(to, from actor.Sender, message string) {
		println(message)
	})

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

		room.Send(text, "")
	}
}

// Topic creates a topic name from an actor name
func topic(name string) string {
	return fmt.Sprintf("actor/%s/", name)
}
