// Package string implements github.com/the-anna-project/clg.Service and
// provides a method to pass through the received string argument as it is. The
// idea behind this pass through CLG is that it does not calculate or manipulate
// any of the given arguments, but only returns them. The idea is that the
// arguments being returned as they are received can be used to be forwarded in
// a specific way based on the forwarding rules associated with a specific
// behaviour ID. This can especially be useful in combination with the buffer
// CLG.
package string

import (
	"sync"

	"github.com/the-anna-project/context"
	"github.com/the-anna-project/id"
)

// ServiceConfig represents the configuration used to create a new CLG service.
type ServiceConfig struct {
	// Dependencies.
	IDService id.Service
}

// DefaultServiceConfig provides a default configuration to create a new CLG
// service by best effort.
func DefaultServiceConfig() ServiceConfig {
	var err error

	var idService id.Service
	{
		idConfig := id.DefaultServiceConfig()
		idService, err = id.NewService(idConfig)
		if err != nil {
			panic(err)
		}
	}

	config := ServiceConfig{
		// Dependencies.
		IDService: idService,
	}

	return config
}

// NewService creates a new configured CLG service.
func NewService(config ServiceConfig) (*Service, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}

	ID, err := config.IDService.New()
	if err != nil {
		return nil, maskAny(err)
	}

	newService := &Service{
		// Internals.
		bootOnce: sync.Once{},
		closer:   make(chan struct{}, 1),
		metadata: map[string]string{
			"id":   ID,
			"kind": "pass/through/string",
			"name": "clg",
			"type": "service",
		},
		shutdownOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Internals.
	bootOnce     sync.Once
	closer       chan struct{}
	metadata     map[string]string
	shutdownOnce sync.Once
}

func (s *Service) Action() interface{} {
	return func(ctx context.Context, s string) string {
		return s
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
