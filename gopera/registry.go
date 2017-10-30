package gopera

type registryInterface interface {
	Configure(...string) error
	RegisterActor(string, OptionsInterface) error
	Watch(func(string, *ActorOptions)) error
	UnregisterActor(string) error
	ParseConfiguration() (map[string]OptionsInterface, error)
}
