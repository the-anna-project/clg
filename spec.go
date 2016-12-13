// Package clg provides the specification of a CLG service.
package clg

// Service represents the CLGs that are interacting with each other within the
// neural network. Each CLG is registered in the neural network. From there
// signals are dispatched across queues in a dynamic fashion until some useful
// calculation took place.
type Service interface {
	// Action returns the CLG's calculation function which implements the CLG's
	// actual business logic.
	Action() interface{}
	// Boot initializes and starts the whole CLG like booting a machine. The call
	// to Boot blocks until the CLG is completely initialized, so you might want
	// to call it in a separate goroutine.
	Boot()
	// Metadata returns a copy of the CLG's metadata. Metadata can be information
	// like service name, service kind, service ID or the like.
	Metadata() map[string]string
	// Shutdown ends all processes of the CLG like shutting down a machine. The
	// call to Shutdown blocks until the CLG is completely shut down, so you might
	// want to call it in a separate goroutine.
	Shutdown()
}
