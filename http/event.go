package http

import (
	"context"
	palikka "palikka-go"
)

type subscriber struct {
	Messages chan []byte
	ctx      context.Context
}

func newSubscriber(ctx context.Context) *subscriber {
	return &subscriber{
		Messages: make(chan []byte),
		ctx:      ctx,
	}
}

type eventService struct {
	SessionStore palikka.SessionStore
	Subscribers  map[string]*subscriber
}

func (s *eventService) Subscribe(id string, subscriber *subscriber) {
	s.Subscribers[id] = subscriber // todo what if same subscriber?
	for {
		select {
		case <-subscriber.ctx.Done():

		}
		}
	}
}

func (s *eventService) Publish() {
	go func() {
		for {
			select {
			case id := <-s.SessionStore.OnDelete():
				s.Subscribers[id].ctx.Done() // todo close somehow
				delete(s.Subscribers, id) // todo close connection
			}
		}
	}()
}

func (s *eventService) delete() {

}