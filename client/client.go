package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type SendArgs struct {
	User string
	Text string
}

type Message struct {
	User string
	Text string
	Time time.Time
}

type HistoryReply struct {
	Messages []Message
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	userRaw, _ := reader.ReadString('\n')
	user := strings.TrimSpace(userRaw)
	if user == "" {
		user = "anonymous"
	}

	var client *rpc.Client
	var err error

	for {
		client, err = rpc.Dial("tcp", "127.0.0.1:1234")
		if err != nil {
			fmt.Println("Could not connect to server, retrying in 2s...")
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	defer client.Close()

	fmt.Println("Connected. Type messages and press Enter. Type 'exit' to quit.")
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		if text == "exit" {
			fmt.Println("Goodbye!")
			return
		}

		args := SendArgs{User: user, Text: text}
		var reply HistoryReply
		callErr := client.Call("Chat.SendMessage", args, &reply)
		if callErr != nil {
			fmt.Println("Error sending message:", callErr)
			fmt.Println("Attempting to reconnect...")
			client, err = rpc.Dial("tcp", "127.0.0.1:1234")
			if err != nil {
				fmt.Println("Reconnect failed, will retry in 2s.")
				time.Sleep(2 * time.Second)
				continue
			}
			fmt.Println("Reconnected. Resending message...")
			callErr = client.Call("Chat.SendMessage", args, &reply)
			if callErr != nil {
				fmt.Println("Failed again:", callErr)
				continue
			}
		}

		// print short history
		fmt.Println("---- Chat History ----")
		for _, m := range reply.Messages {
			fmt.Printf("[%s] %s: %s\n", m.Time.Format("15:04:05"), m.User, m.Text)
		}
		fmt.Println("----------------------")
	}
}
