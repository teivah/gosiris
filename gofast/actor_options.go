package gofast

type ActorOptions struct {
	remote          bool
	serverAlias     string
	connectionAlias string
	endpoint        string
}

type OptionsInterface interface {
	SetRemote(bool) OptionsInterface
	SetServerAlias(string) OptionsInterface
	SetConnectionAlias(string) OptionsInterface
	SetEndpoint(string) OptionsInterface
	Remote() bool
	ServerAlias() string
	ConnectionAlias() string
	Endpoint() string
}

func (options *ActorOptions) SetRemote(b bool) OptionsInterface {
	options.remote = b
	return options
}

func (options *ActorOptions) SetServerAlias(s string) OptionsInterface {
	options.serverAlias = s
	return options
}

func (options *ActorOptions) SetConnectionAlias(s string) OptionsInterface {
	options.connectionAlias = s
	return options
}

func (options *ActorOptions) SetEndpoint(s string) OptionsInterface {
	options.endpoint = s
	return options
}

func (options *ActorOptions) Remote() bool {
	return options.remote
}

func (options *ActorOptions) ServerAlias() string {
	return options.serverAlias
}

func (options *ActorOptions) ConnectionAlias() string {
	return options.connectionAlias
}

func (options *ActorOptions) Endpoint() string {
	return options.endpoint
}
