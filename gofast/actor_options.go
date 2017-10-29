package gofast

type ActorOptions struct {
	local           bool
	connectionAlias string
	endpoint        string
}

type OptionsInterface interface {
	Local(bool) OptionsInterface
	ConnectionAlias(string) OptionsInterface
	Endpoint(string) OptionsInterface
}

func (options *ActorOptions) Local(b bool) OptionsInterface {
	options.local = b
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
