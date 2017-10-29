package gofast

type ActorOptions struct {
	local       bool
	serverAlias string
	endpoint    string
}

type OptionsInterface interface {
	Local(bool) OptionsInterface
	ServerAlias(string) OptionsInterface
	Endpoint(string) OptionsInterface
}

func (options *ActorOptions) Local(b bool) OptionsInterface {
	options.local = b
	return options
}

func (options *ActorOptions) ServerAlias(s string) OptionsInterface {
	options.serverAlias = s
	return options
}

func (options *ActorOptions) Endpoint(s string) OptionsInterface {
	options.endpoint = s
	return options
}
