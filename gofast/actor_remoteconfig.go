package gofast

var conf map[string]actorRemoteConfiguration

type actorRemoteConfiguration struct {
	remoteType  string
	url         string
	destination string
}

type actorRemoteConfigurationInterface interface {
	Configure(...string) error
	ParseConfiguration() (map[string]actorRemoteConfiguration, error)
}

func ActorRemoteConfiguration() map[string]actorRemoteConfiguration {
	return conf
}
