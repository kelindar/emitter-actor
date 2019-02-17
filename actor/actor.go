package actor

import (
	"fmt"
	"strings"

	emitter "github.com/emitter-io/go/v2"
)

// Key for 'actor/#/' channel, allowing us to do everything with any actor.
//const key = "LkmxQlmzyBsByET468R1AWxh79Qu1kMS"

// Handler represents a message handler
type handler = func(to, from Sender, message string)

// Sender represents an object which can receive messages.
type Sender interface {
	Send(command, message string)
}

// Topic creates a topic name from an actor name
func topic(name string) string {
	return fmt.Sprintf("actor/%s/", name)
}

// Command creates a command with a message
func command(from, command, message string) string {
	return fmt.Sprintf("%s %s %s", from, command, message)
}

// remote represents a remote actor
type remote struct {
	send func(cmd, message string) error
}

// Remote creates a way of talking to a remote actor
func Remote(key, to, from string, network *emitter.Client) Sender {
	return &remote{send: func(cmd, message string) error {
		return network.Publish(key, topic(to), command(from, cmd, message))
	}}
}

// Send sends a message to the actor.
func (r remote) Send(cmd, message string) {
	_ = r.send(cmd, message)
}

// Actor represents a game actor
type Actor struct {
	Sender
	name     string
	pubkey   string
	handlers map[string]handler
	network  *emitter.Client
}

// New creates a new actor
func New(subKey, pubKey, name string, network *emitter.Client, private bool) (actor *Actor, err error) {
	actor = &Actor{
		Sender:   Remote(pubKey, name, name, network),
		handlers: make(map[string]handler),
		pubkey:   pubKey,
		name:     name,
		network:  network,
	}

	// Create a private link so we can receive dedicated replies
	topic := topic(name)
	if private {
		link, err := network.CreatePrivateLink(subKey, "actor/", "s", false)
		if err != nil {
			return nil, err
		}

		topic = link.Channel
	}

	// Subscribe to the channel
	err = network.Subscribe(subKey, topic, actor.onMessageReceived)
	return
}

// Occurs when a remote message is received
func (a *Actor) onMessageReceived(_ *emitter.Client, msg emitter.Message) {

	// Get the command and message from the payload
	request := strings.SplitN(string(msg.Payload()), " ", 3)
	if len(request) != 3 {
		return
	}

	from := request[0]
	command := request[1]
	message := request[2]
	if fn, ok := a.handlers[command]; ok {
		fn(a, Remote(a.pubkey, from, a.name, a.network), message)
	}
}

// On attaches a message handler.
func (a *Actor) On(command string, fn handler) {
	a.handlers[command] = fn
}
