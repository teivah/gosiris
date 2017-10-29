package gofast

//type actorRemoteConfiguration struct {
//	remoteType  string
//	url         string
//	destination string
//}

type actorRemoteConfigurationInterface interface {
	Configure(...string) error
	RegisterActor(string, OptionsInterface) error
	UnregisterActor(string) error
	ParseConfiguration() (map[string]OptionsInterface, error)
	Configuration() map[string]OptionsInterface
}
