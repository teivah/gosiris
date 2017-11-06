package gosiris

type registryInterface interface {
	Configure(...string) error
	Close()
	RegisterActor(string, OptionsInterface) error
	Watch(func(string, *ActorOptions), func(string)) error
	UnregisterActor(string) error
	ParseConfiguration() (map[string]OptionsInterface, error)
}
