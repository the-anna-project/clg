// Package isgreater implements github.com/the-anna-project/clg.Service and
// provides a method to identify if the first given number is greater than the
// later.
package isgreater

import (
	"sync"

	"github.com/the-anna-project/context"
	"github.com/the-anna-project/id"

	"github.com/the-anna-project/clg"
)

// Config represents the configuration used to create a new CLG service.
type Config struct {
	// Dependencies.
	IDService id.Service
}

// DefaultConfig provides a default configuration to create a new CLG service by
// best effort.
func DefaultConfig() Config {
	var err error

	var idService id.Service
	{
		idConfig := id.DefaultConfig()
		idService, err = id.New(idConfig)
		if err != nil {
			panic(err)
		}
	}

	config := Config{
		// Dependencies.
		IDService: idService,
	}

	return config
}

// New creates a new configured CLG service.
func New(config Config) (clg.Service, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}

	ID, err := config.IDService.New()
	if err != nil {
		panic(err)
	}

	newService := &service{
		// Internals.
		bootOnce: sync.Once{},
		closer:   make(chan struct{}, 1),
		metadata: map[string]string{
			"id":   ID,
			"kind": "isgreater",
			"name": "clg",
			"type": "service",
		},
		shutdownOnce: sync.Once{},
	}

	return newService, nil
}

type service struct {
	// Internals.
	bootOnce     sync.Once
	closer       chan struct{}
	metadata     map[string]string
	shutdownOnce sync.Once
}

func (s *service) Action() interface{} {
	return func(ctx context.Context, a, b float64) bool {
		return a > b
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
