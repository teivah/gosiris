package gopera

func RootActor() *Actor {
	return root
}

type Actor struct {
	name     string
	conf     map[string]func(Message)
	dataChan chan Message
	parent   actorInterface
	unbecome map[string]func(Message)
}

type RemoteActor struct {
	Actor
}

type actorInterface interface {
	React(string, func(Message)) *Actor
	reactions() map[string]func(Message)
	unbecomeHistory() map[string]func(Message)
	getDataChan() chan Message
	setDataChan(chan Message)
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

func (actor *Actor) React(messageType string, f func(Message)) *Actor {
	if actor.conf == nil {
		actor.conf = make(map[string]func(Message))
		actor.unbecome = make(map[string]func(Message))
	}

	actor.conf[messageType] = f

	return actor
}

func (actor *Actor) reactions() map[string]func(Message) {
	return actor.conf
}

func (actor *Actor) unbecomeHistory() map[string]func(Message) {
	return actor.unbecome
}

func (actor *Actor) getDataChan() chan Message {
	return actor.dataChan
}

func (actor *Actor) setDataChan(dataChan chan Message) {
	actor.dataChan = dataChan
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
