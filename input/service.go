// Package input implements github.com/the-anna-project/clg.Service and provides
// the entry to the neural network. When being executed the CLGs action tries to
// lookup the information peer associated with the given information sequence.
// In case the information peer cannot be found within the connection space, a
// new information peer is created. In any case the ID of the information peer
// is added to the given context. Further CLGs may or may not make use of it.
package input

import (
	"sync"

	"github.com/the-anna-project/context"
	informationid "github.com/the-anna-project/context/information/id"
	"github.com/the-anna-project/id"
	"github.com/the-anna-project/peer"

	"github.com/the-anna-project/clg"
)

// Config represents the configuration used to create a new CLG service.
type Config struct {
	// Dependencies.
	IDService      id.Service
	PeerCollection *peer.Collection
}

// DefaultConfig provides a default configuration to create a new CLG service by
// best effort.
func DefaultConfig() Config {
	var err error

	var idService id.Service
	{
		idConfig := id.DefaultServiceConfig()
		idService, err = id.NewService(idConfig)
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

	config := Config{
		// Dependencies.
		IDService:      idService,
		PeerCollection: peerCollection,
	}

	return config
}

// New creates a new configured CLG service.
func New(config Config) (clg.Service, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if config.PeerCollection == nil {
		return nil, maskAnyf(invalidConfigError, "peer collection must not be empty")
	}

	ID, err := config.IDService.New()
	if err != nil {
		return nil, maskAny(err)
	}

	newService := &service{
		// Internals.
		bootOnce: sync.Once{},
		closer:   make(chan struct{}, 1),
		metadata: map[string]string{
			"id":   ID,
			"kind": "input",
			"name": "clg",
			"type": "service",
		},
		shutdownOnce: sync.Once{},

		// Settings.
		peer: config.PeerCollection,
	}

	return newService, nil
}

type service struct {
	// Internals.
	bootOnce     sync.Once
	closer       chan struct{}
	metadata     map[string]string
	shutdownOnce sync.Once

	// Dependencies.
	peer *peer.Collection
}

func (s *service) Action() interface{} {
	return func(ctx context.Context, informationSequence string) error {
		informationPeer, err := s.peer.Information.Search(informationSequence)
		if peer.IsNotFound(err) {
			// The given information sequence was never seen before. Thus we register
			// it now by creating an information peer for it.
			informationPeer, err = s.peer.Information.Create(informationSequence)
			if err != nil {
				return maskAny(err)
			}
		} else if err != nil {
			return maskAny(err)
		}

		ctx = informationid.NewContext(ctx, informationPeer.ID())

		return nil
	}
}

func (s *service) Boot() {
	s.bootOnce.Do(func() {
		// Service specific boot logic goes here.
	})
}

func (s *service) Metadata() map[string]string {
	m := map[string]string{}
	for k, v := range s.metadata {
		m[k] = v
	}
	return m
}

func (s *service) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.closer)
	})
}
