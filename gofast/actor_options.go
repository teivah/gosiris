package gofast

type ActorOptions struct {
	serverAlias     string
	connectionAlias string
	endpoint        string
}

type OptionsInterface interface {
	SetServerAlias(string) OptionsInterface
	SetConnectionAlias(string) OptionsInterface
	SetEndpoint(string) OptionsInterface
	ServerAlias() string
	ConnectionAlias() string
	Endpoint() string
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

func (options *ActorOptions) ServerAlias() string {
	return options.serverAlias
}

func (options *ActorOptions) ConnectionAlias() string {
	return options.connectionAlias
}

func (options *ActorOptions) Endpoint() string {
	return options.endpoint
}
