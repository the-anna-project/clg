// Package input implements github.com/the-anna-project/clg.Service and provides
// the entry to the neural network. When being executed the CLGs action fetches
// the information ID associated with the given information sequence. In case
// the information sequence is not found within the underlying storage, a new
// information ID is generated and used to store the given information sequence.
// In any case the information ID is added to the given context.
package input

import (
	"fmt"
	"sync"

	"github.com/the-anna-project/context"
	informationid "github.com/the-anna-project/context/information/id"
	"github.com/the-anna-project/id"
	storagecollection "github.com/the-anna-project/storage/collection"
	storageerror "github.com/the-anna-project/storage/error"

	"github.com/the-anna-project/clg"
)

// Config represents the configuration used to create a new CLG service.
type Config struct {
	// Dependencies.
	IDService         id.Service
	StorageCollection *storagecollection.Collection
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

	var storageCollection *storagecollection.Collection
	{
		storageConfig := storagecollection.DefaultConfig()
		storageCollection, err = storagecollection.New(storageConfig)
		if err != nil {
			panic(err)
		}
	}

	config := Config{
		// Dependencies.
		IDService:         idService,
		StorageCollection: storageCollection,
	}

	return config
}

// New creates a new configured CLG service.
func New(config Config) (clg.Service, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if config.StorageCollection == nil {
		return nil, maskAnyf(invalidConfigError, "storage collection must not be empty")
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
		idService: config.IDService,
		storage:   config.StorageCollection,
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
	idService id.Service
	storage   *storagecollection.Collection
}

func (s *service) Action() interface{} {
	return func(ctx context.Context, informationSequence string) error {
		// TODO there should be a service to manage information sequences and
		// information IDs etc.
		informationIDKey := fmt.Sprintf("information-sequence:%s:information-id", informationSequence)
		informationID, err := s.storage.General.Get(informationIDKey)
		if storageerror.IsNotFound(err) {
			// The given information sequence was never seen before. Thus we register it
			// now with its own very unique information ID.
			newID, err := s.idService.New()
			if err != nil {
				return maskAny(err)
			}
			informationID = string(newID)

			err = s.storage.General.Set(informationIDKey, informationID)
			if err != nil {
				return maskAny(err)
			}

			informationSequenceKey := fmt.Sprintf("information-id:%s:information-sequence", informationID)
			err = s.storage.General.Set(informationSequenceKey, informationSequence)
			if err != nil {
				return maskAny(err)
			}
		} else if err != nil {
			return maskAny(err)
		}

		ctx = informationid.NewContext(ctx, informationID)

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
