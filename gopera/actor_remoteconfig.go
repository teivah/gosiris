package gopera

type actorRemoteConfigurationInterface interface {
	Configure(...string) error
	RegisterActor(string, OptionsInterface) error
	UnregisterActor(string) error
	ParseConfiguration() (map[string]OptionsInterface, error)
	Configuration() map[string]OptionsInterface
}
