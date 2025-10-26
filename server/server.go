package main

import (
	"errors"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type Message struct {
	User string
	Text string
	Time time.Time
}

type SendArgs struct {
	User string
	Text string
}

type HistoryReply struct {
	Messages []Message
}

type Chat struct {
	mu       sync.Mutex
	messages []Message
}

func (c *Chat) SendMessage(args SendArgs, reply *HistoryReply) error {
	if args.Text == "" {
		return errors.New("empty message")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	msg := Message{
		User: args.User,
		Text: args.Text,
		Time: time.Now(),
	}
	c.messages = append(c.messages, msg)
	// return copy of history
	h := make([]Message, len(c.messages))
	copy(h, c.messages)
	reply.Messages = h
	return nil
}

func (c *Chat) GetHistory(_ struct{}, reply *HistoryReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	h := make([]Message, len(c.messages))
	copy(h, c.messages)
	reply.Messages = h
	return nil
}

func main() {
	chat := new(Chat)
	err := rpc.Register(chat)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", ":1234") // port 1234
	if err != nil {
		panic(err)
	}
	defer l.Close()
	println("Chat server listening on :1234")
	for {
		conn, err := l.Accept()
		if err != nil {
			println("accept error:", err.Error())
			continue
		}
		go rpc.ServeConn(conn)
	}
}
