// Package output implements github.com/the-anna-project/clg.Service and
// provides one of the two very special CLGs. That is the output CLG. Its
// purpose is to check if the calculated output matches the provided
// expectation, if any expectation given. The output CLG is handled in a special
// way because it determines the end of all requested calculations within the
// neural network. After the output CLG has been executed, the calculated output
// is returned back to the requesting client.
package output

import (
	"reflect"
	"sync"

	"github.com/the-anna-project/context"
	currentbehaviourid "github.com/the-anna-project/context/current/behaviour/id"
	destinationid "github.com/the-anna-project/context/destination/id"
	"github.com/the-anna-project/context/expectation"
	firstbehaviourid "github.com/the-anna-project/context/first/behaviour/id"
	firstinformationid "github.com/the-anna-project/context/first/information/id"
	sourceids "github.com/the-anna-project/context/source/ids"
	"github.com/the-anna-project/event"
	"github.com/the-anna-project/id"
	"github.com/the-anna-project/output"
	"github.com/the-anna-project/peer"
)

// ServiceConfig represents the configuration used to create a new CLG service.
type ServiceConfig struct {
	// Dependencies.
	EventCollection  *event.Collection
	IDService        id.Service
	OutputCollection *output.Collection
	PeerCollection   *peer.Collection
}

// DefaultServiceConfig provides a default configuration to create a new CLG
// service by best effort.
func DefaultServiceConfig() ServiceConfig {
	var err error

	var eventCollection *event.Collection
	{
		eventConfig := event.DefaultCollectionConfig()
		eventCollection, err = event.NewCollection(eventConfig)
		if err != nil {
			panic(err)
		}
	}

	var idService id.Service
	{
		idConfig := id.DefaultServiceConfig()
		idService, err = id.NewService(idConfig)
		if err != nil {
			panic(err)
		}
	}

	var outputCollection *output.Collection
	{
		outputConfig := output.DefaultCollectionConfig()
		outputCollection, err = output.NewCollection(outputConfig)
		if err != nil {
			panic(err)
		}
	}

	var peerCollection *peer.Collection
	{
		peerConfig := peer.DefaultCollectionConfig()
		peerCollection, err = peer.NewCollection(peerConfig)
		if err != nil {
			panic(err)
		}
	}

	config := ServiceConfig{
		// Dependencies.
		EventCollection:  eventCollection,
		IDService:        idService,
		OutputCollection: outputCollection,
		PeerCollection:   peerCollection,
	}

	return config
}

// NewService creates a new configured CLG service.
func NewService(config ServiceConfig) (*Service, error) {
	// Dependencies.
	if config.EventCollection == nil {
		return nil, maskAnyf(invalidConfigError, "event collection must not be empty")
	}
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if config.OutputCollection == nil {
		return nil, maskAnyf(invalidConfigError, "output collection must not be empty")
	}
	if config.PeerCollection == nil {
		return nil, maskAnyf(invalidConfigError, "peer collection must not be empty")
	}

	ID, err := config.IDService.New()
	if err != nil {
		return nil, maskAny(err)
	}

	newService := &Service{
		// Dependencies.
		event:  config.EventCollection,
		output: config.OutputCollection,
		peer:   config.PeerCollection,

		// Internals.
		bootOnce: sync.Once{},
		closer:   make(chan struct{}, 1),
		metadata: map[string]string{
			"id":   ID,
			"kind": "output",
			"name": "clg",
			"type": "service",
		},
		shutdownOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	event  *event.Collection
	output *output.Collection
	peer   *peer.Collection

	// Internals.
	bootOnce     sync.Once
	closer       chan struct{}
	metadata     map[string]string
	shutdownOnce sync.Once
}

func (s *Service) Action() interface{} {
	return func(ctx context.Context, informationSequence string) error {
		// Check the calculated output against the provided expectation, if any. In
		// case there is no expectation provided, we simply go with what we
		// calculated. This then means we are probably not in a training situation.
		e, ok := expectation.FromContext(ctx)
		if !ok {
			err := s.sendTextOutput(ctx, informationSequence)
			if err != nil {
				return maskAny(err)
			}

			return nil
		}

		// There is an expectation provided. Thus we are going to check the calculated
		// output against it. In case the provided expectation does match the
		// calculated result, we simply return it.
		calculatedOutput := e.Output()
		if informationSequence == calculatedOutput {
			err := s.sendTextOutput(ctx, informationSequence)
			if err != nil {
				return maskAny(err)
			}
		}

		// The calculated output did not match the given expectation. That means we
		// need to calculate some new output to match the given expectation. To do so
		// we create a new network payload and assign the input CLG of the current CLG
		// tree to it by queueing the new network payload in the underlying storage.
		err := s.forwardNetworkPayload(ctx)
		if err != nil {
			return maskAny(err)
		}

		// The calculated output did not match the given expectation. We return an
		// error to let the neural network know about it.
		return maskAnyf(expectationNotMetError, "'%s' != '%s'", informationSequence, calculatedOutput)
	}
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		// Service specific boot logic goes here.
	})
}

func (s *Service) Metadata() map[string]string {
	m := map[string]string{}
	for k, v := range s.metadata {
		m[k] = v
	}
	return m
}

func (s *Service) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.closer)
	})
}

// TODO there is no CLG to read from the certenty pyramid

func (s *Service) forwardNetworkPayload(ctx context.Context) error {
	// When the output CLG forwards signals to the network, it forwards to the
	// input CLG. Therefore we want to find the very first information sequence
	// provided by the client. The information sequence is obtained by an
	// information peer, which is located somewhere within the connection space.
	firstInformationID, ok := firstinformationid.FromContext(ctx)
	if !ok {
		return maskAnyf(invalidInformationIDError, "must not be empty")
	}
	firstInformationPeer, err := s.peer.Information.SearchByID(firstInformationID)
	if err != nil {
		return maskAny(err)
	}

	// We want to forward the current signal to the input CLG. Therefore we have
	// to find the first ID of the very first behaviour of the current CLG tree.
	firstBehaviourID, ok := firstbehaviourid.FromContext(ctx)
	if !ok {
		return maskAnyf(invalidBehaviourIDError, "must not be empty")
	}

	// To be able to forward the current signal to the first behaviour of the
	// current CLG tree, we have to reference the current CLG as sender in the
	// queue event.
	currentBehaviourID, ok := currentbehaviourid.FromContext(ctx)
	if !ok {
		return maskAnyf(invalidBehaviourIDError, "must not be empty")
	}

	// Here we set the required control flow information. The first behaviour ID
	// is the behaviour ID of the input CLG within the current CLG tree. The
	// current behaviour ID is the behaviour ID of the current output CLG within
	// the current CLG tree.
	ctx = destinationid.NewContext(ctx, firstBehaviourID)
	ctx = sourceids.NewContext(ctx, []string{currentBehaviourID})

	// Now the signal will be created which we will publish to the signal queue
	// below. The signal has necessary information applied to keep the neural
	// networks moving.
	signalConfig := event.DefaultSignalConfig()
	signalConfig.Arguments = []reflect.Value{reflect.ValueOf(firstInformationPeer.Value())}
	signalConfig.Context = ctx
	signal, err := event.NewSignal(signalConfig)
	if err != nil {
		return maskAny(err)
	}

	// Finally, publish the created signal. Somewhere som worker will pick up this
	// specific event to process it. Then a new calculation iteration begins.
	err = s.event.Signal.Publish(signal)
	if err != nil {
		return maskAny(err)
	}

	return nil
}

func (s *Service) sendTextOutput(ctx context.Context, informationSequence string) error {
	config := output.DefaultConfig()
	config.Text = informationSequence
	newOutput, err := output.New(config)
	if err != nil {
		return maskAny(err)
	}

	s.output.Text.Channel() <- newOutput

	return nil
}
