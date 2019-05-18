package adapters

import (
	"fmt"
	"log"
	"strings"
)

// DrivingAdapter defines a driving adapter interface.
// A driving adapter is an adapter that turns concrete protocol inputs to commands
type DrivingAdapter interface {
	Describe() DrivingAdapterDescription
	Start()
	Stop()
	GetCommand() DrivingAdapterCmd
}

// DrivingAdapterCmd defines a generic command input
type DrivingAdapterCmd struct {
	CmdHandler func()
}

// DrivingAdapterDescription defines a description structure for driving adapters
type DrivingAdapterDescription struct {
	Title string
	State DrivingAdapterState
	// Capabilities (list of supported commands)
}

// DrivingAdapterConfig defines a generic driving adapter configuration
type DrivingAdapterConfig map[string]string

// DrivingAdapterState defines various states for driving adapters
type DrivingAdapterState int

// Constants for driving adapters' states
const (
	Unknown DrivingAdapterState = iota
	Initialized
	Started
	Stopped
)

// DrivingAdapterFactory factory for instantiating driving adapters
type DrivingAdapterFactory func(config DrivingAdapterConfig) (DrivingAdapter, error)

var drivingAdapterFactories = make(map[string]DrivingAdapterFactory)

// RegisterDrivingAdapter registers a driving adapter
func RegisterDrivingAdapter(name string, factory DrivingAdapterFactory) {
	if factory == nil {
		log.Panicf("DrivingAdapter factory %s does not exist.", name)
	}
	_, registered := drivingAdapterFactories[name]
	if registered {
		log.Printf("Datastore factory %s already registered. Ignoring.", name)
	}
	drivingAdapterFactories[name] = factory
}

// CreateDrivingAdapter creates a driving adapter by calling its factory method
func CreateDrivingAdapter(conf DrivingAdapterConfig) (DrivingAdapter, error) {

	engineName := conf["adapter-type"]

	engineFactory, ok := drivingAdapterFactories[engineName]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging
		availableAdapters := make([]string, len(drivingAdapterFactories))
		for f := range drivingAdapterFactories {
			availableAdapters = append(availableAdapters, f)
		}
		return nil, fmt.Errorf("Invalid DrivenAdapter name. Must be one of: %s", strings.Join(availableAdapters, ", "))
	}

	// Run the factory with the configuration
	return engineFactory(conf)
}
