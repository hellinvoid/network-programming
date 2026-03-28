package service

import (
	"errors"
	"fmt"
	"strings"
)

type ChatRoom map[string]chan string

func NewChatRoom() ChatRoom {
	return make(ChatRoom)
}

func (cr ChatRoom) Join(username string) (chan string, chan string, error) {

	// Check if username already exists
	if _, exists := cr[username]; exists {
		return nil, nil, errors.New("Username already exists")
	}
	// Make a new send channel per user
	cr[username] = make(chan string)

	// Send join message to eveyone
	joinMessage := fmt.Sprintf("* %s has entered the room\n", username)
	go cr.Publish(username, joinMessage)

	// Send members currently in room to joined user
	go cr.SendMembersInRoom(username)

	// Make a receive channel for chat room to receive messages from user
	crReceive := make(chan string)
	go cr.ListenForMessages(username, crReceive)

	// send of chat room = receive of user
	// receive of chat room = send of user
	// Return receive of user, send of user
	return cr[username], crReceive, nil
}

func (cr ChatRoom) Leave(username string) {
	// Remove the user from chat room
	delete(cr, username)

	// Publish leave message
	leaveMessage := fmt.Sprintf("* %s has left the room\n", username)
	go cr.Publish(username, leaveMessage)
}

func (cr ChatRoom) Publish(usernameToSkip, message string) {
	// Send message to all present in room except the user himself
	for key, ch := range cr {
		if key == usernameToSkip {
			continue
		}
		ch <- message
	}
}

func (cr ChatRoom) SendMembersInRoom(username string) {
	members := make([]string, 0)

	for key := range cr {
		if key == username {
			continue
		}
		members = append(members, key)
	}

	membersMessage := fmt.Sprintf("* The room contains: %s\n", strings.Join(members, ", "))
	cr[username] <- membersMessage
}

func (cr ChatRoom) ListenForMessages(username string, crReceive chan string) {
	defer cr.Leave(username)

	for msg := range crReceive {
		cr.Publish(username, msg)
	}
}
