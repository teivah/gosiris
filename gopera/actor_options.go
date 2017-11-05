package gopera

const (
	defaultBufferSize = 64
)

type ActorOptions struct {
	parent      string
	remote      bool //Default: false
	autoclose   bool //Default: false
	remoteType  string
	url         string
	destination string
	bufferSize  int //Default: 64
}

type OptionsInterface interface {
	SetUrl(string) OptionsInterface
	SetRemote(bool) OptionsInterface
	SetAutoclose(bool) OptionsInterface
	Remote() bool
	Autoclose() bool
	SetRemoteType(string) OptionsInterface
	RemoteType() string
	Url() string
	SetDestination(string) OptionsInterface
	Destination() string
	setParent(string)
	Parent() string
	SetBufferSize(int)
	BufferSize() int
}

func (options *ActorOptions) SetRemote(b bool) OptionsInterface {
	options.remote = b
	return options
}

func (options *ActorOptions) Remote() bool {
	return options.remote
}

func (options *ActorOptions) SetAutoclose(b bool) OptionsInterface {
	options.autoclose = b
	return options
}

func (options *ActorOptions) Autoclose() bool {
	return options.autoclose
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

func (options *ActorOptions) setParent(s string) {
	options.parent = s
}

func (options *ActorOptions) Parent() string {
	return options.parent
}

func (options *ActorOptions) SetBufferSize(i int) {
	options.bufferSize = i
}

func (options *ActorOptions) BufferSize() int {
	return options.bufferSize
}
