package clg

import (
	"sync"

	"github.com/the-anna-project/id"
	"github.com/the-anna-project/peer"

	"github.com/the-anna-project/clg/divide"
	"github.com/the-anna-project/clg/greater"
	"github.com/the-anna-project/clg/input"
	"github.com/the-anna-project/clg/isbetween"
	"github.com/the-anna-project/clg/isgreater"
	"github.com/the-anna-project/clg/islesser"
	"github.com/the-anna-project/clg/lesser"
	"github.com/the-anna-project/clg/multiply"
)

// CollectionConfig represents the configuration used to create a new CLG
// collection.
type CollectionConfig struct {
	// Dependencies.
	IDService      id.Service
	PeerCollection *peer.Collection
}

// DefaultCollectionConfig provides a default configuration to create a new CLG
// collection by best effort.
func DefaultCollectionConfig() CollectionConfig {
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

	config := CollectionConfig{
		// Dependencies.
		IDService:      idService,
		PeerCollection: peerCollection,
	}

	return config
}

// NewCollection creates a new configured CLG Collection.
func NewCollection(config CollectionConfig) (*Collection, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if config.PeerCollection == nil {
		return nil, maskAnyf(invalidConfigError, "peer collection must not be empty")
	}

	var err error

	var divideService Service
	{
		divideConfig := divide.DefaultServiceConfig()
		divideConfig.IDService = config.IDService
		divideService, err = divide.NewService(divideConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var greaterService Service
	{
		greaterConfig := greater.DefaultServiceConfig()
		greaterConfig.IDService = config.IDService
		greaterService, err = greater.NewService(greaterConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var inputService Service
	{
		inputConfig := input.DefaultServiceConfig()
		inputConfig.IDService = config.IDService
		inputConfig.PeerCollection = config.PeerCollection
		inputService, err = input.NewService(inputConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isBetweenService Service
	{
		isBetweenConfig := isbetween.DefaultServiceConfig()
		isBetweenConfig.IDService = config.IDService
		isBetweenService, err = isbetween.NewService(isBetweenConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isGreaterService Service
	{
		isGreaterConfig := isgreater.DefaultServiceConfig()
		isGreaterConfig.IDService = config.IDService
		isGreaterService, err = isgreater.NewService(isGreaterConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isLesserService Service
	{
		isLesserConfig := islesser.DefaultServiceConfig()
		isLesserConfig.IDService = config.IDService
		isLesserService, err = islesser.NewService(isLesserConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var lesserService Service
	{
		lesserConfig := lesser.DefaultServiceConfig()
		lesserConfig.IDService = config.IDService
		lesserService, err = lesser.NewService(lesserConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var multiplyService Service
	{
		multiplyConfig := multiply.DefaultServiceConfig()
		multiplyConfig.IDService = config.IDService
		multiplyService, err = multiply.NewService(multiplyConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	newCollection := &Collection{
		// Internals.
		bootOnce:     sync.Once{},
		shutdownOnce: sync.Once{},

		// Public.
		List: []Service{
			divideService,
			greaterService,
			inputService,
			isBetweenService,
			isGreaterService,
			isLesserService,
			lesserService,
			multiplyService,
		},

		Divide:    divideService,
		Greater:   greaterService,
		Input:     inputService,
		IsBetween: isBetweenService,
		IsGreater: isGreaterService,
		IsLesser:  isLesserService,
		Lesser:    lesserService,
		Multiply:  multiplyService,
	}

	return newCollection, nil
}

// Collection is the object bundling all CLGs.
type Collection struct {
	// Internals.
	bootOnce     sync.Once
	shutdownOnce sync.Once

	// Public.
	List []Service

	Divide    Service
	Greater   Service
	Input     Service
	IsBetween Service
	IsGreater Service
	IsLesser  Service
	Lesser    Service
	Multiply  Service
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
