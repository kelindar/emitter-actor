package main

import (
	emitter "github.com/emitter-io/go/v2"
	"github.com/kelindar/emitter-actor/actor"
)

// Key for 'actor/#/' channel, allowing us to do everything with any actor.
const key = "LkmxQlmzyBsByET468R1AWxh79Qu1kMS"

func main() {
	client, err := emitter.Connect("", func(c *emitter.Client, msg emitter.Message) {
		println("unknown message: ", msg.Topic(), string(msg.Payload()))
	})
	if err != nil {
		panic(err)
	}

	// First room
	room1, _ := actor.New(key, key, "room1", client, false)
	room1.On("enter", func(to, from actor.Sender, message string) {
		from.Send("tell", "You've entered a dark room.")
	})
	room1.On("look", func(to, from actor.Sender, message string) {
		from.Send("tell", "You notice a small lamp on the desk.")
	})
	room1.On("lamp", func(to, from actor.Sender, message string) {
		from.Send("tell", "You turn on the lamp and now you can see a door.")
	})
	room1.On("door", func(to, from actor.Sender, message string) {
		from.Send("move", "room2")
	})

	// Second room
	room2, _ := actor.New(key, key, "room2", client, false)
	room2.On("enter", func(to, from actor.Sender, message string) {
		from.Send("tell", "You've entered a big room.")
	})
	room2.On("look", func(to, from actor.Sender, message string) {
		from.Send("tell", "The lights are on and you can see a door and a poster.")
	})
	room2.On("door", func(to, from actor.Sender, message string) {
		from.Send("move", "room1")
	})
	room2.On("poster", func(to, from actor.Sender, message string) {
		from.Send("tell", "The poster says 'visit our Github to know more about emitter!'")
	})

	println("server started")
	for {
	}
}
