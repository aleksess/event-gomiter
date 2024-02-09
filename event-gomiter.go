package eventgomiter

type EventIdentifier string

type Event interface{}

type EventHandler func(Event)

type EventHandleable interface {
	EventHandler() (EventHandler, error)
}

type EventEmitter interface {
	Emit(EventIdentifier, Event) error
	RegisterHandler(EventIdentifier, EventHandler) error
	RegisterHandleable(EventIdentifier, EventHandleable) error
}

type ChannelEventEmitter struct {
	channels map[EventIdentifier]chan Event
	handlers map[EventIdentifier][]EventHandler
}

var _ EventEmitter = (*ChannelEventEmitter)(nil)

func NewChannelEventEmitter() *ChannelEventEmitter {
	emitter := &ChannelEventEmitter{
		channels: make(map[EventIdentifier]chan Event),
		handlers: make(map[EventIdentifier][]EventHandler),
	}

	go func() {
		for {
			for eventId, channel := range emitter.channels {
				event := <-channel
				for _, handler := range emitter.handlers[eventId] {
					go handler(event)
				}
			}
		}
	}()

	return emitter
}

func (c *ChannelEventEmitter) Emit(id EventIdentifier, event Event) error {
	c.channels[id] <- event
	return nil
}

func (c *ChannelEventEmitter) RegisterHandler(id EventIdentifier, handler EventHandler) error {
	if c.handlers[id] == nil {
		c.handlers[id] = make([]EventHandler, 0, 1)
	}
	if c.channels[id] == nil {
		c.channels[id] = make(chan Event)
	}

	c.handlers[id] = append(c.handlers[id], handler)

	return nil
}

func (c *ChannelEventEmitter) RegisterHandleable(id EventIdentifier, handleable EventHandleable) error {
	handler, err := handleable.EventHandler()
	if err != nil {
		return err
	} else {
		return c.RegisterHandler(id, handler)
	}
}
