package gopera

type ActorOptions struct {
	remote      bool
	remoteType  string
	url         string
	destination string
}

type OptionsInterface interface {
	SetRemote(bool) OptionsInterface
	Remote() bool
	SetRemoteType(string) OptionsInterface
	RemoteType() string
	SetUrl(string) OptionsInterface
	Url() string
	SetDestination(string) OptionsInterface
	Destination() string
}

func (options *ActorOptions) SetRemote(b bool) OptionsInterface {
	options.remote = b
	return options
}

func (options *ActorOptions) Remote() bool {
	return options.remote
}

func (options *ActorOptions) SetRemoteType(s string) OptionsInterface {
	options.remoteType = s
	return options
}

func (options *ActorOptions) RemoteType() string {
	return options.remoteType
}

func (options *ActorOptions) SetUrl(s string) OptionsInterface {
	options.url = s
	return options
}

func (options *ActorOptions) Url() string {
	return options.url
}

func (options *ActorOptions) SetDestination(s string) OptionsInterface {
	options.destination = s
	return options
}

func (options *ActorOptions) Destination() string {
	return options.destination
}
