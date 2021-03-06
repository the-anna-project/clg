package clg

import (
	"sync"

	"github.com/the-anna-project/event"
	"github.com/the-anna-project/id"
	"github.com/the-anna-project/index"
	"github.com/the-anna-project/output"
	"github.com/the-anna-project/peer"
	"github.com/the-anna-project/random"

	divideclg "github.com/the-anna-project/clg/divide"
	greaterclg "github.com/the-anna-project/clg/greater"
	inputclg "github.com/the-anna-project/clg/input"
	isbetweenclg "github.com/the-anna-project/clg/is/between"
	isgreaterclg "github.com/the-anna-project/clg/is/greater"
	islesserclg "github.com/the-anna-project/clg/is/lesser"
	lesserclg "github.com/the-anna-project/clg/lesser"
	multiplyclg "github.com/the-anna-project/clg/multiply"
	outputclg "github.com/the-anna-project/clg/output"
	passthroughfloat64clg "github.com/the-anna-project/clg/pass/through/float64"
	passthroughstringclg "github.com/the-anna-project/clg/pass/through/string"
	readinformationsequence "github.com/the-anna-project/clg/read/information/sequence"
	readseparatorclg "github.com/the-anna-project/clg/read/separator"
	roundclg "github.com/the-anna-project/clg/round"
	subtractclg "github.com/the-anna-project/clg/subtract"
	sumclg "github.com/the-anna-project/clg/sum"
)

// CollectionConfig represents the configuration used to create a new CLG
// collection.
type CollectionConfig struct {
	// Dependencies.
	EventCollection  *event.Collection
	IDService        id.Service
	IndexService     index.Service
	OutputCollection *output.Collection
	PeerCollection   *peer.Collection
	RandomService    random.Service
}

// DefaultCollectionConfig provides a default configuration to create a new CLG
// collection by best effort.
func DefaultCollectionConfig() CollectionConfig {
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

	var indexService index.Service
	{
		indexConfig := index.DefaultServiceConfig()
		indexService, err = index.NewService(indexConfig)
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

	var randomService random.Service
	{
		randomConfig := random.DefaultServiceConfig()
		randomService, err = random.NewService(randomConfig)
		if err != nil {
			panic(err)
		}
	}

	config := CollectionConfig{
		// Dependencies.
		EventCollection:  eventCollection,
		IDService:        idService,
		IndexService:     indexService,
		OutputCollection: outputCollection,
		PeerCollection:   peerCollection,
		RandomService:    randomService,
	}

	return config
}

// NewCollection creates a new configured CLG Collection.
func NewCollection(config CollectionConfig) (*Collection, error) {
	// Dependencies.
	if config.EventCollection == nil {
		return nil, maskAnyf(invalidConfigError, "event collection must not be empty")
	}
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if config.IndexService == nil {
		return nil, maskAnyf(invalidConfigError, "index service must not be empty")
	}
	if config.OutputCollection == nil {
		return nil, maskAnyf(invalidConfigError, "output collection must not be empty")
	}
	if config.PeerCollection == nil {
		return nil, maskAnyf(invalidConfigError, "peer collection must not be empty")
	}
	if config.RandomService == nil {
		return nil, maskAnyf(invalidConfigError, "random service must not be empty")
	}

	var err error

	var divideService Service
	{
		divideConfig := divideclg.DefaultServiceConfig()
		divideConfig.IDService = config.IDService
		divideService, err = divideclg.NewService(divideConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var greaterService Service
	{
		greaterConfig := greaterclg.DefaultServiceConfig()
		greaterConfig.IDService = config.IDService
		greaterService, err = greaterclg.NewService(greaterConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var inputService Service
	{
		inputConfig := inputclg.DefaultServiceConfig()
		inputConfig.IDService = config.IDService
		inputConfig.PeerCollection = config.PeerCollection
		inputService, err = inputclg.NewService(inputConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isBetweenService Service
	{
		isBetweenConfig := isbetweenclg.DefaultServiceConfig()
		isBetweenConfig.IDService = config.IDService
		isBetweenService, err = isbetweenclg.NewService(isBetweenConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isGreaterService Service
	{
		isGreaterConfig := isgreaterclg.DefaultServiceConfig()
		isGreaterConfig.IDService = config.IDService
		isGreaterService, err = isgreaterclg.NewService(isGreaterConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var isLesserService Service
	{
		isLesserConfig := islesserclg.DefaultServiceConfig()
		isLesserConfig.IDService = config.IDService
		isLesserService, err = islesserclg.NewService(isLesserConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var lesserService Service
	{
		lesserConfig := lesserclg.DefaultServiceConfig()
		lesserConfig.IDService = config.IDService
		lesserService, err = lesserclg.NewService(lesserConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var multiplyService Service
	{
		multiplyConfig := multiplyclg.DefaultServiceConfig()
		multiplyConfig.IDService = config.IDService
		multiplyService, err = multiplyclg.NewService(multiplyConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var outputService Service
	{
		outputConfig := outputclg.DefaultServiceConfig()
		outputConfig.EventCollection = config.EventCollection
		outputConfig.IDService = config.IDService
		outputConfig.OutputCollection = config.OutputCollection
		outputConfig.PeerCollection = config.PeerCollection
		outputService, err = outputclg.NewService(outputConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var passThroughFloat64Service Service
	{
		passThroughFloat64Config := passthroughfloat64clg.DefaultServiceConfig()
		passThroughFloat64Config.IDService = config.IDService
		passThroughFloat64Service, err = passthroughfloat64clg.NewService(passThroughFloat64Config)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var passThroughStringService Service
	{
		passThroughStringConfig := passthroughstringclg.DefaultServiceConfig()
		passThroughStringConfig.IDService = config.IDService
		passThroughStringService, err = passthroughstringclg.NewService(passThroughStringConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var readInformationSequenceService Service
	{
		readInformationSequenceConfig := readinformationsequence.DefaultServiceConfig()
		readInformationSequenceConfig.IDService = config.IDService
		readInformationSequenceConfig.PeerCollection = config.PeerCollection
		readInformationSequenceService, err = readinformationsequence.NewService(readInformationSequenceConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var readSeparatorService Service
	{
		readSeparatorConfig := readseparatorclg.DefaultServiceConfig()
		readSeparatorConfig.IDService = config.IDService
		readSeparatorConfig.IndexService = config.IndexService
		readSeparatorConfig.PeerCollection = config.PeerCollection
		readSeparatorConfig.RandomService = config.RandomService
		readSeparatorService, err = readseparatorclg.NewService(readSeparatorConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var roundService Service
	{
		roundConfig := roundclg.DefaultServiceConfig()
		roundConfig.IDService = config.IDService
		roundService, err = roundclg.NewService(roundConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var subtractService Service
	{
		subtractConfig := subtractclg.DefaultServiceConfig()
		subtractConfig.IDService = config.IDService
		subtractService, err = subtractclg.NewService(subtractConfig)
		if err != nil {
			return nil, maskAny(err)
		}
	}

	var sumService Service
	{
		sumConfig := sumclg.DefaultServiceConfig()
		sumConfig.IDService = config.IDService
		sumService, err = sumclg.NewService(sumConfig)
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
			outputService,
			passThroughFloat64Service,
			passThroughStringService,
			readInformationSequenceService,
			readSeparatorService,
			roundService,
			subtractService,
			sumService,
		},

		Divide:                  divideService,
		Greater:                 greaterService,
		Input:                   inputService,
		IsBetween:               isBetweenService,
		IsGreater:               isGreaterService,
		IsLesser:                isLesserService,
		Lesser:                  lesserService,
		Multiply:                multiplyService,
		Output:                  outputService,
		PassThroughFloat64:      passThroughFloat64Service,
		PassThroughString:       passThroughStringService,
		ReadInformationSequence: readInformationSequenceService,
		ReadSeparator:           readSeparatorService,
		Round:                   roundService,
		Subtract:                subtractService,
		Sum:                     sumService,
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

	Divide                  Service
	Greater                 Service
	Input                   Service
	IsBetween               Service
	IsGreater               Service
	IsLesser                Service
	Lesser                  Service
	Multiply                Service
	Output                  Service
	PassThroughFloat64      Service
	PassThroughString       Service
	ReadInformationSequence Service
	ReadSeparator           Service
	Round                   Service
	Subtract                Service
	Sum                     Service
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
