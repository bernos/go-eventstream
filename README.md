# go-eventstream

A package for creating composable concurrent event streams. 

The go programming language features primitives such as channels and go routines for writing concurrent programs. While these primitives can be used to implement powerful concurrency patterns, doing so over and over again can be tedious and error-prone.

The go-eventstream package provides a small set of simple, composable abstractions over these base primitives. Inspired by concepts from functional reactive programming, these abstractions can be used to compose powerful, concurrent programs that are easy to reason about, and easy to test.

## Events, Streams and Transformers

The Event interface represents a single event that will be sent on an event stream. Each event contains either a value, or an error.

```go
type Event interface {
	Value() interface{}
	Error() error
}
```

The Stream interface represents a continuous stream of `Event`s. Event streams are a wrapper arround `chan Event` and can be created from many sources, from a `[]interface{}`, to a polling `func()(interface{},error)`, to an existing channel.

```go
type Stream interface {
  	// Retrieve the underlying Event channel from the stream
	Events() <-chan Event
	
	// Send a value and/or error on the stream
	Send(interface{}, error)
	
	// Cancel the stream. In the case of a root stream (ie. one with no parent), 
	// this will close the underlying Event channel. For streams that were created via
	// CreateChild, Cancel will defer to calling Cancel on the the parent stream
	Cancel()
	
	// Create a child stream. The new stream will have its own underlying
	// Event channel, but its Cancel method will cascade back to the root
	// of the stream heirarchy, closing the Event channel of the first 
	// ancestor Stream.
	CreateChild(chan Event) Stream

  	// ...
}
```

The `Transformer` interface represents the ability to transform one `Stream` to another. `Transformer` can be used to implement concepts such as filtering, mapping, folding and throttling. `Transformer` can also compose with itself, making it easy to build powerful concurrent pipelines with ease. The go-eventstream package ships with implementations of many common Stream Transformers, such as `Map`, `PMap`, `FlatMap`, `PFlatMap`, `Scan`, `Reduce` and more.

```go
type Transformer interface {
	Transform(Stream) Stream
	Compose(Transformer) Transformer
	
	// ...
}
```

