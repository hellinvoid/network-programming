package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/hellinvoid/network-programming/problems/003-budget-chat/utils"
)

const (
	MAX_USERNAME_LENGTH = 16
)

type Client struct {
	username string
	conn     net.Conn
	receive  chan string
	send     chan string
}

func NewClient(conn net.Conn, cr ChatRoom) (*Client, error) {

	// Asking for name
	greetMessage := "Welcome to budgetchat! What shall I call you?\n"
	_, err := io.WriteString(conn, greetMessage)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(conn)
	usernameBytes, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	username := string(usernameBytes)
	
	// Check for invalid username 
	if len(username) > MAX_USERNAME_LENGTH || len(username) < 1 || !utils.IsAlnum(username) {
		return nil, errors.New("Username invalid")
	}

	// Initialize receive and send channels
	receive, send, err := cr.Join(username)
	if err != nil {
		return nil, err
	}


	return &Client{
		username: username,
		conn:     conn,
		receive:  receive,
		send:     send,
	}, nil
}

func (cl *Client) SendMessageToRoom() error {
	r := bufio.NewReader(cl.conn)
	messageBytes, _, err := r.ReadLine()
	if err != nil {
		return err
	}

	// Format and send the message from user to send channel
	message := string(messageBytes)
	formattedMessage := fmt.Sprintf("[%s] %s\n", cl.username, message)
	cl.send <- formattedMessage
	return nil
}

func (cl *Client) ReceiveMessagesFromRoom() {
	// Receive messages from receive channel and send to user via TCP
	for message := range cl.receive {
		_, err := io.WriteString(cl.conn, message)
		if err != nil {
			return
		}
	}
}

func (cl *Client) Start() {
	// Close send and receive channels at end
	defer close(cl.send)
	defer close(cl.receive)

	go cl.ReceiveMessagesFromRoom()
	for {
		err := cl.SendMessageToRoom()
		if err != nil {
			return
		}
	}

}
