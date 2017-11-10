package gosiris

import "github.com/opentracing/opentracing-go"

func RootActor() *Actor {
	return root
}

type Actor struct {
	name      string
	conf      map[string]func(Context)
	dataChan  chan Context
	closeChan chan interface{}
	parent    actorInterface
	unbecome  map[string]func(Context)
	span      opentracing.Span
}

type RemoteActor struct {
	Actor
}

type actorInterface interface {
	React(string, func(Context)) *Actor
	reactions() map[string]func(Context)
	unbecomeHistory() map[string]func(Context)
	getDataChan() chan Context
	setDataChan(chan Context)
	getCloseChan() chan interface{}
	setCloseChan(chan interface{})
	setName(string)
	setParent(actorInterface)
	Parent() ActorRefInterface
	Name() string
	Close()
}

func (actor *Actor) Close() {
	ActorSystem().closeLocalActor(actor.name)
}

func (actor *Actor) String() string {
	return actor.name
}

func (actor *Actor) React(messageType string, f func(Context)) *Actor {
	if actor.conf == nil {
		actor.conf = make(map[string]func(Context))
		actor.unbecome = make(map[string]func(Context))
	}

	actor.conf[messageType] = f

	return actor
}

func (actor *Actor) reactions() map[string]func(Context) {
	return actor.conf
}

func (actor *Actor) unbecomeHistory() map[string]func(Context) {
	return actor.unbecome
}

func (actor *Actor) getDataChan() chan Context {
	return actor.dataChan
}

func (actor *Actor) setDataChan(dataChan chan Context) {
	actor.dataChan = dataChan
}

func (actor *Actor) getCloseChan() chan interface{} {
	return actor.closeChan
}

func (actor *Actor) setCloseChan(closeChan chan interface{}) {
	actor.closeChan = closeChan
}

func (actor *Actor) setName(name string) {
	actor.name = name
}

func (actor *Actor) setParent(parent actorInterface) {
	actor.parent = parent
}

func (actor *Actor) Parent() ActorRefInterface {
	parent, _ := ActorSystem().ActorOf(actor.parent.Name())
	return parent
}

func (actor *Actor) Name() string {
	return actor.name
}
