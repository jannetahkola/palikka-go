package http

import (
	"context"
	"fmt"
	"github.com/coder/websocket"
	"net/http"
	palikka "palikka-go"
	"sync"
	"time"
)

func (s *Server) handleUpgradeConn(w http.ResponseWriter, r *http.Request) {
	session := palikka.SessionFromContext(r.Context())

	var lock sync.Mutex
	var c *websocket.Conn
	var closed bool

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("accept error: %v", err))
		return
	}

	lock.Lock()
	if closed {
		lock.Unlock()
		return
	}

	lock.Unlock()
	defer c.CloseNow()

	ctx := c.CloseRead(context.Background())
	sub := newSubscriber(ctx)
	s.EventService.Subscribe(session.ID, sub)

	for {
		select {
		case msg := <-sub.Messages:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				//return err
			}
		case id := <-s.SessionStore.OnDelete():
			// todo remove from event sub map
			if id == session.ID {
				c.CloseNow()
			}
		case <-ctx.Done():
			// todo what is this really
			//return ctx.Err()
		}
	}

	//session.Values["Ws-Conn"] = c
	//session.OnDelete = func(s *palikka.Session) {
	//	conn := s.Values["Ws-Conn"].(*websocket.Conn)
	//	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
	//		fmt.Println(fmt.Errorf("error closing websocket connection %w", err))
	//	}
	//}
}
