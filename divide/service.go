// Package divide implements github.com/the-anna-project/clg.Service and
// provides the mathematical operation of division.
package divide

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
			"kind": "divide",
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
	return func(ctx context.Context, a, b float64) float64 {
		return a / b
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
