// Package collection provides a bundle of CLG services.
package collection

import (
	"sync"

	"github.com/the-anna-project/id"

	"github.com/the-anna-project/clg"
	"github.com/the-anna-project/clg/divide"
	"github.com/the-anna-project/clg/greater"
)

// Config represents the configuration used to create a new CLG collection.
type Config struct {
	// Dependencies.
	IDService id.Service
}

// DefaultConfig provides a default configuration to create a new CLG collection
// by best effort.
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

// New creates a new configured CLG Collection.
func New(config Config) (*Collection, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}

	var err error

	var divideService clg.Service
	{
		divideConfig := divide.DefaultConfig()
		divideConfig.IDService = config.IDService
		divideService, err = divide.New(divideConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var greaterService clg.Service
	{
		greaterConfig := greater.DefaultConfig()
		greaterConfig.IDService = config.IDService
		greaterService, err = greater.New(greaterConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	newCollection := &Collection{
		// Internals.
		bootOnce:     sync.Once{},
		shutdownOnce: sync.Once{},

		// Public.
		Divide:  divideService,
		Greater: greaterService,
	}

	return newCollection, nil
}

// Collection is the object bundling all CLGs.
type Collection struct {
	// Internals.
	bootOnce     sync.Once
	shutdownOnce sync.Once

	// Public.
	Divide  clg.Service
	Greater clg.Service
}

func (c *Collection) Boot() {
	c.bootOnce.Do(func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			c.Divide.Boot()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			c.Greater.Boot()
			wg.Done()
		}()

		wg.Wait()
	})
}

func (c *Collection) Shutdown() {
	c.shutdownOnce.Do(func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			c.Divide.Shutdown()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			c.Greater.Shutdown()
			wg.Done()
		}()

		wg.Wait()
	})
}
