package service

import (
	"context"
	"io"
)

// EventHandler is a callback function, which is executed when a subscriber
// receives a message. Note that this function does not return an error, because
// the message bus does not know how to handle that error and would simply
// cancel the subscription. Errors have to be handled inside the event handler.
type EventHandler func(ctx context.Context, msg []byte)

// MessageBus defines the interface for publishing messages to a topic and
// subscribing for receiving messages from a topic.
type MessageBus interface {
	io.Closer

	// Publish publishes a message to a given topic. This
	// function returns [ErrConnectionClosed] in case the
	// connection to the message broker is closed.
	Publish(_ context.Context, topic string, msg []byte) error

	// Subscribe subscribes to the given topic. The event handler
	// callback will be executed on every received message. This
	// function returns [ErrConnectionClosed] in case the
	// connection to the message broker is closed. This is a
	// blocking function. Canceling the context will cancel the
	// subscription.
	Subscribe(_ context.Context, topic string, h EventHandler) error
}
