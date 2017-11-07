package gosiris

import "time"

const (
	defaultBufferSize = 64
)

type ActorOptions struct {
	parent         string
	remote         bool //Default: false
	autoclose      bool //Default: false
	remoteType     string
	url            string
	destination    string
	bufferSize     int //Default: 64
	defaultWatcher time.Duration
}

//TODO No interface
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
	SetBufferSize(int) OptionsInterface
	BufferSize() int
	SetDefaultWatcher(time.Duration) OptionsInterface
	DefaultWatcher() time.Duration
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

func (options *ActorOptions) SetBufferSize(i int) OptionsInterface {
	options.bufferSize = i
	return options
}

func (options *ActorOptions) BufferSize() int {
	return options.bufferSize
}

func (options *ActorOptions) SetDefaultWatcher(d time.Duration) OptionsInterface {
	options.defaultWatcher = d
	return options
}

func (options *ActorOptions) DefaultWatcher() time.Duration {
	return options.defaultWatcher
}
