// Package collection provides a bundle of CLG services.
package collection

import (
	"sync"

	"github.com/the-anna-project/id"

	"github.com/the-anna-project/clg"
	"github.com/the-anna-project/clg/divide"
	"github.com/the-anna-project/clg/greater"
	"github.com/the-anna-project/clg/input"
	"github.com/the-anna-project/clg/isbetween"
	"github.com/the-anna-project/clg/isgreater"
	"github.com/the-anna-project/clg/islesser"
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

	var inputService clg.Service
	{
		inputConfig := input.DefaultConfig()
		inputConfig.IDService = config.IDService
		inputService, err = input.New(inputConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isBetweenService clg.Service
	{
		isBetweenConfig := isbetween.DefaultConfig()
		isBetweenConfig.IDService = config.IDService
		isBetweenService, err = isbetween.New(isBetweenConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isGreaterService clg.Service
	{
		isGreaterConfig := isgreater.DefaultConfig()
		isGreaterConfig.IDService = config.IDService
		isGreaterService, err = isgreater.New(isGreaterConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isLesserService clg.Service
	{
		isLesserConfig := islesser.DefaultConfig()
		isLesserConfig.IDService = config.IDService
		isLesserService, err = islesser.New(isLesserConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	newCollection := &Collection{
		// Internals.
		bootOnce:     sync.Once{},
		shutdownOnce: sync.Once{},

		// Public.
		List: []clg.Service{
			divideService,
			greaterService,
			inputService,
			isBetweenService,
			isGreaterService,
			isLesserService,
		},

		Divide:    divideService,
		Greater:   greaterService,
		Input:     inputService,
		IsBetween: isBetweenService,
		IsGreater: isGreaterService,
		IsLesser:  isLesserService,
	}

	return newCollection, nil
}

// Collection is the object bundling all CLGs.
type Collection struct {
	// Internals.
	bootOnce     sync.Once
	shutdownOnce sync.Once

	// Public.
	List []clg.Service

	Divide    clg.Service
	Greater   clg.Service
	Input     clg.Service
	IsBetween clg.Service
	IsGreater clg.Service
	IsLesser  clg.Service
}

func (c *Collection) Boot() {
	c.bootOnce.Do(func() {
		var wg sync.WaitGroup

		for _, s := range c.List {
			wg.Add(1)
			go func() {
				s.Boot()
				wg.Done()
			}()
		}

		wg.Wait()
	})
}

func (c *Collection) Shutdown() {
	c.shutdownOnce.Do(func() {
		var wg sync.WaitGroup

		for _, s := range c.List {
			wg.Add(1)
			go func() {
				s.Shutdown()
				wg.Done()
			}()
		}

		wg.Wait()
	})
}
