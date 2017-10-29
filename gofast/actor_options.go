package gofast

type ActorOptions struct {
	serverAlias     string
	connectionAlias string
	endpoint        string
}

type OptionsInterface interface {
	ServerAlias(string) OptionsInterface
	ConnectionAlias(string) OptionsInterface
	Endpoint(string) OptionsInterface
}

func (options *ActorOptions) ServerAlias(s string) OptionsInterface {
	options.serverAlias = s
	return options
}

func (options *ActorOptions) ConnectionAlias(s string) OptionsInterface {
	options.connectionAlias = s
	return options
}

func (options *ActorOptions) Endpoint(s string) OptionsInterface {
	options.endpoint = s
	return options
}
