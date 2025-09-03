package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	channelPkg "github.com/khabib-developer/hydra-server/pkg/channel"
	"github.com/khabib-developer/hydra-server/pkg/user"
	"github.com/khabib-developer/hydra-server/pkg/websocket"
)

type Channel struct {
	ID    uuid.UUID
	Name  string
	Owner *User
	Users []*User
}

func (server *Server) createChannel(payload json.RawMessage, sender *User) {
	var channelName string
	if err := json.Unmarshal(payload, &channelName); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Unsupported type of channel name")
		return
	}

	for _, ch := range server.Channels {
		if ch.Name == channelName {
			server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Channel already exists")
			return
		}
	}

	if sender.JoinedChannel != nil {
		for i, u := range sender.JoinedChannel.Users {
			if u == sender {
				sender.JoinedChannel.Users = append(sender.JoinedChannel.Users[:i], sender.JoinedChannel.Users[i+1:]...)
				break
			}
		}
		sender.JoinedChannel = nil
	}

	channel := &Channel{
		Name:  channelName,
		Owner: sender,
		Users: []*User{sender},
	}

	sender.JoinedChannel = channel

	server.Channels = append(server.Channels, channel)

	server.sendRawMessage(sender.SafeConn, websocket.MessageTypeCreate, channel.Name)
}

func (server *Server) joinChannel(payload json.RawMessage, sender *User) {
	var channelName string
	if err := json.Unmarshal(payload, &channelName); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Unsupperted type of channel name")
		return
	}

	var channelToJoin *Channel
	for _, ch := range server.Channels {
		if ch.Name == channelName {
			channelToJoin = ch
			break
		}
	}

	if channelToJoin == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Channel not found")
		return
	}

	if sender.JoinedChannel != nil {
		for i, u := range sender.JoinedChannel.Users {
			if u == sender {
				fmt.Println("Removing user from previous channel")
				sender.JoinedChannel.Users = append(sender.JoinedChannel.Users[:i], sender.JoinedChannel.Users[i+1:]...)
				break
			}
		}
		sender.JoinedChannel = nil
	}

	sender.JoinedChannel = channelToJoin
	channelToJoin.Users = append(channelToJoin.Users, sender)
	server.sendRawMessage(sender.SafeConn, websocket.MessageTypeJoin, channelToJoin.Name)
}

func (server *Server) broadcastMessage(payload json.RawMessage, sender *User) {
	var message string
	if err := json.Unmarshal(payload, &message); err != nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Invalid message payload")
		return
	}

	if sender.JoinedChannel == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "You are not in a channel")
		return
	}

	channelMessage := channelPkg.ChannelMessageDto{
		Channel: sender.JoinedChannel.Name,
		Sender:  sender.Username,
		Message: message,
	}

	jsonMessage, err := json.Marshal(channelMessage)
	if err != nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Failed to marshal channel message")
		return
	}

	for _, user := range sender.JoinedChannel.Users {
		server.sendMessage(user.SafeConn, websocket.MessageTypeBroadcast, jsonMessage)
	}
}

func (server *Server) destroyChannel(payload json.RawMessage, sender *User) {
	var channelName string

	if err := json.Unmarshal(payload, &channelName); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Unsupported type of channel name")
		return
	}

	var channelToDestroy *Channel
	var index int = -1

	for i, ch := range server.Channels {
		if channelName == ch.Name {
			channelToDestroy = ch
			index = i
			break
		}
	}

	if channelToDestroy == nil {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "Channel not found")
		return
	}

	if channelToDestroy.Owner != sender {
		server.sendRawMessage(sender.SafeConn, websocket.MessageTypeError, "You are not the owner of this channel")
		return
	}

	for _, user := range channelToDestroy.Users {
		user.JoinedChannel = nil
		server.sendRawMessage(user.SafeConn, websocket.MessageTypeInfo, fmt.Sprintf("Channel '%s' has been destroyed", channelToDestroy.Name))
	}

	if index != -1 {
		server.Channels = append(server.Channels[:index], server.Channels[index+1:]...)
	}
}

func (server *Server) destroyUserChannels(u *User) {
	for _, channel := range server.Channels {
		if channel.Owner == u {

			channelNameJson, err := json.Marshal(channel.Name)

			if err != nil {
				return
			}

			server.destroyChannel(channelNameJson, u)

		}
	}
}

func (server *Server) GetChannels(w http.ResponseWriter, r *http.Request) {
	channels := make([]channelPkg.ChannelDto, len(server.Channels))
	i := 0
	connID := r.Header.Get("connID")
	user := server.Users[connID]

	for _, channel := range server.Channels {
		channels[i] = channelPkg.ChannelDto{
			Name:    channel.Name,
			Members: len(channel.Users),
			Mine:    user == channel.Owner,
		}
		i++
	}
	msg, err := json.Marshal(channels)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
func (server *Server) GetChannelMembers(w http.ResponseWriter, r *http.Request) {
	users := make([]user.UserDTO, 0, len(server.Users))
	for _, userItem := range server.Users {
		users = append(users, user.UserDTO{
			Username: userItem.Username,
			Private:  len(userItem.Password) > 0,
		})
	}

	msg, err := json.Marshal(users)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
