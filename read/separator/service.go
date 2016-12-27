// Package separator implements github.com/the-anna-project/clg.Service and
// provides functionality to provide a separator stored as peer value of a
// specific information peer. When this CLG is being executed, it uses the
// context to identify itself. The context contains information about the CLGs
// behaviour ID, which is used to lookup a mapping pointing to an information
// peer ID. This information peer ID is used to lookup the actual information
// peer and its associated value, which is the separator. In case there is no
// mapping for the current behaviour ID, a new separator will be made up and a
// new information peer as well as the necessary index mapping. In any case a
// separator will be returned.
package separator

import (
	"sync"

	"github.com/the-anna-project/context"
	currentbehaviourid "github.com/the-anna-project/context/current/behaviour/id"
	"github.com/the-anna-project/id"
	"github.com/the-anna-project/index"
	"github.com/the-anna-project/peer"
	"github.com/the-anna-project/random"
)

const (
	// NamespaceBehaviourID represents the namespace being used to map a specific
	// behaviour ID to a specific information ID using the index service.
	NamespaceBehaviourID = "behaviour-id"
	// NamespaceInformationID represents the namespace being used to map a
	// specific behaviour ID to a specific information ID using the index service.
	NamespaceInformationID = "information-id"
	// NamespaceSeparator represents the namespace being used to map a specific
	// behaviour ID to a specific information ID using the index service.
	NamespaceSeparator = "separator"
)

// ServiceConfig represents the configuration used to create a new CLG service.
type ServiceConfig struct {
	// Dependencies.
	IDService      id.Service
	IndexService   index.Service
	PeerCollection *peer.Collection
	RandomService  random.Service
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

	var indexService index.Service
	{
		indexConfig := index.DefaultServiceConfig()
		indexService, err = index.NewService(indexConfig)
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

	config := ServiceConfig{
		// Dependencies.
		IDService:      idService,
		IndexService:   indexService,
		PeerCollection: peerCollection,
		RandomService:  randomService,
	}

	return config
}

// NewService creates a new configured CLG service.
func NewService(config ServiceConfig) (*Service, error) {
	// Dependencies.
	if config.IDService == nil {
		return nil, maskAnyf(invalidConfigError, "ID service must not be empty")
	}
	if config.IndexService == nil {
		return nil, maskAnyf(invalidConfigError, "index service must not be empty")
	}
	if config.PeerCollection == nil {
		return nil, maskAnyf(invalidConfigError, "peer collection must not be empty")
	}
	if config.RandomService == nil {
		return nil, maskAnyf(invalidConfigError, "random service must not be empty")
	}

	ID, err := config.IDService.New()
	if err != nil {
		return nil, maskAny(err)
	}

	newService := &Service{
		// Dependencies.
		index:  config.IndexService,
		peer:   config.PeerCollection,
		random: config.RandomService,

		// Internals.
		bootOnce: sync.Once{},
		closer:   make(chan struct{}, 1),
		metadata: map[string]string{
			"id":   ID,
			"kind": "read/separator",
			"name": "clg",
			"type": "service",
		},
		shutdownOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	index  index.Service
	peer   *peer.Collection
	random random.Service

	// Internals.
	bootOnce     sync.Once
	closer       chan struct{}
	metadata     map[string]string
	shutdownOnce sync.Once
}

func (s *Service) Action() interface{} {
	return func(ctx context.Context) (string, error) {
		behaviourID, ok := currentbehaviourid.FromContext(ctx)
		if !ok {
			return "", maskAnyf(invalidBehaviourIDError, "must not be empty")
		}

		informationID, err := s.index.Search(NamespaceSeparator, NamespaceBehaviourID, NamespaceInformationID, behaviourID)
		if index.IsNotFound(err) {
			// Create a new random separator. Therefore we lookup some random
			// information peer and use its value for the new separator.
			//
			// TODO we use one single character of the information peer value as
			// separator for now. There might be applications in the future benefiting
			// from separators having multiple characters.
			informationPeer, err := s.peer.Information.Random()
			if err != nil {
				return "", maskAny(err)
			}
			feature := informationPeer.Value()
			featureIndex, err := s.random.CreateMax(len(feature))
			if err != nil {
				return "", maskAny(err)
			}
			separator := string(feature[featureIndex])

			// Create a new information peer and the necessary mapping so we can lookup
			// the separator when the current CLG is executed again using its very
			// unique behaviour ID.
			informationPeer, err = s.peer.Information.Create(separator)
			if err != nil {
				return "", maskAny(err)
			}
			err = s.index.Create(NamespaceSeparator, NamespaceBehaviourID, NamespaceInformationID, behaviourID, informationPeer.ID())
			if err != nil {
				return "", maskAny(err)
			}

			// We created the information peer for the new separator and the necessary
			// index mapping between the current behaviour ID and the new information
			// ID. We are done and can savely return the new separator.
			return separator, nil
		} else if err != nil {
			return "", maskAny(err)
		}

		// We found an information ID using an existing index mapping between the
		// current behaviour ID and its associated information ID. We lookup the peer
		// and return the separator obtained by the information peer.
		informationPeer, err := s.peer.Information.SearchByID(informationID)
		if err != nil {
			return "", maskAny(err)
		}

		return informationPeer.Value(), nil
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
